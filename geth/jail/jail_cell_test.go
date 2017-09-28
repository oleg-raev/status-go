package jail

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/robertkrimen/otto"
	"github.com/stretchr/testify/require"
)

func TestJailLoopInCall(t *testing.T) {
	vm := otto.New()
	cell, err := newCell("some-id", vm)
	require.NoError(t, err)

	items := make(chan string)

	err = cell.Set("__captureResponse", func(val string) otto.Value {
		go func() { items <- val }()
		return otto.UndefinedValue()
	})
	require.NoError(t, err)

	_, err = cell.Run(`
		function callRunner(namespace){
			return setTimeout(function(){
				__captureResponse(namespace);
			}, 1000);
		}
	`)
	require.NoError(t, err)

	_, err = cell.Call("callRunner", nil, "softball")
	require.NoError(t, err)

	select {
	case received := <-items:
		require.Equal(t, received, "softball")
	case <-time.After(5 * time.Second):
		require.FailNow(t, "Failed to received event response")
	}
}

// TestJailLoopRace tests multiple setTimeout callbacks,
// supposed to be run with '-race' flag.
func TestJailLoopRace(t *testing.T) {
	vm := otto.New()
	cell, err := newCell("some-id", vm)
	require.NoError(t, err)

	items := make(chan struct{})

	err = cell.Set("__captureResponse", func() otto.Value {
		go func() { items <- struct{}{} }()
		return otto.UndefinedValue()
	})
	require.NoError(t, err)

	_, err = cell.Run(`
		function callRunner(){
			return setTimeout(function(){
				__captureResponse();
			}, 1000);
		}
	`)
	require.NoError(t, err)

	for i := 0; i < 100; i++ {
		_, err = cell.Call("callRunner", nil)
		require.NoError(t, err)
	}

	for i := 0; i < 100; i++ {
		select {
		case <-items:
		case <-time.After(5 * time.Second):
			require.FailNow(t, "test timed out")
		}
	}
}

func TestJailFetchPromise(t *testing.T) {
	body := `{"key": "value"}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(body))
	}))
	defer server.Close()

	vm := otto.New()
	cell, err := newCell("some-id", vm)
	require.NoError(t, err)

	dataCh := make(chan otto.Value, 1)
	errCh := make(chan otto.Value, 1)

	err = cell.Set("__captureSuccess", func(res otto.Value) { dataCh <- res })
	require.NoError(t, err)
	err = cell.Set("__captureError", func(res otto.Value) { errCh <- res })
	require.NoError(t, err)

	// run JS code for fetching valid URL
	_, err = cell.Run(`fetch('` + server.URL + `').then(function(r) {
		return r.text()
	}).then(function(data) {
		__captureSuccess(data)
	}).catch(function (e) {
		__captureError(e)
	})`)
	require.NoError(t, err)

	select {
	case data := <-dataCh:
		require.True(t, data.IsString())
		require.Equal(t, body, data.String())
	case err := <-errCh:
		require.Fail(t, "request failed", err)
	case <-time.After(time.Second):
		require.FailNow(t, "test timed out")
	}
}

func TestJailFetchCatch(t *testing.T) {
	vm := otto.New()
	cell, err := newCell("some-id", vm)
	require.NoError(t, err)

	dataCh := make(chan otto.Value, 1)
	errCh := make(chan otto.Value, 1)

	err = cell.Set("__captureSuccess", func(res otto.Value) { dataCh <- res })
	require.NoError(t, err)
	err = cell.Set("__captureError", func(res otto.Value) { errCh <- res })
	require.NoError(t, err)

	// run JS code for fetching invalid URL
	_, err = cell.Run(`fetch('http://👽/nonexistent').then(function(r) {
		return r.text()
	}).then(function(data) {
		__captureSuccess(data)
	}).catch(function (e) {
		__captureError(e)
	})`)
	require.NoError(t, err)

	select {
	case data := <-dataCh:
		require.Fail(t, "request should have failed, but returned", data)
	case e := <-errCh:
		require.True(t, e.IsObject())
		name, err := e.Object().Get("name")
		require.NoError(t, err)
		require.Equal(t, "Error", name.String())
		_, err = e.Object().Get("message")
		require.NoError(t, err)
	case <-time.After(time.Second):
		require.FailNow(t, "test timed out")
	}
}

// TestJailFetchRace tests multiple fetch callbacks,
// supposed to be run with '-race' flag.
func TestJailFetchRace(t *testing.T) {
	body := `{"key": "value"}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(body))
	}))
	defer server.Close()

	vm := otto.New()
	cell, err := newCell("some-id", vm)
	require.NoError(t, err)

	dataCh := make(chan otto.Value, 1)
	errCh := make(chan otto.Value, 1)

	err = cell.Set("__captureSuccess", func(res otto.Value) { dataCh <- res })
	require.NoError(t, err)
	err = cell.Set("__captureError", func(res otto.Value) { errCh <- res })
	require.NoError(t, err)

	// run JS code for fetching valid URL
	_, err = cell.Run(`fetch('` + server.URL + `').then(function(r) {
		return r.text()
	}).then(function(data) {
		__captureSuccess(data)
	}).catch(function (e) {
		__captureError(e)
	})`)
	require.NoError(t, err)

	// run JS code for fetching invalid URL
	_, err = cell.Run(`fetch('http://👽/nonexistent').then(function(r) {
		return r.text()
	}).then(function(data) {
		__captureSuccess(data)
	}).catch(function (e) {
		__captureError(e)
	})`)
	require.NoError(t, err)

	for i := 0; i < 2; i++ {
		select {
		case data := <-dataCh:
			require.True(t, data.IsString())
			require.Equal(t, body, data.String())
		case e := <-errCh:
			require.True(t, e.IsObject())
			name, err := e.Object().Get("name")
			require.NoError(t, err)
			require.Equal(t, "Error", name.String())
			_, err = e.Object().Get("message")
			require.NoError(t, err)
		case <-time.After(time.Second):
			require.FailNow(t, "test timed out")
		}
	}
}
