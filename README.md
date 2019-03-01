# Proxy
A http proxy in golang. Not currently working for HTTPs but stay tuned...



Specification:

1.Respond to HTTP & HTTPS requests,and should display each request on a management console. It should forward the request to the Web server and relay the response to the browser.

2.Handle websocket connections.

3.Dynamically block selected URLs via the management console.

4.Efficiently cache requests locally and thus save bandwidth. You must gather timing and bandwidth data to prove the efficiency of your proxy.

5.Handle multiple requests simultaneously by implementing a threaded server.
