var socket = io({ transports: ["websocket"] });

function disconnected(message) {
  $("#connection-status").removeClass("connected");
  $("#connection-status").text("Disconnected");
}

function disconnectedFromBackendService(message) {
  $("#connection-status").removeClass("connected");
  $("#connection-status").text("Counting Service is Unreachable");
}

socket.on("disconnect", disconnected);
socket.on("connect_error", disconnected);
socket.on("connect_timeout", disconnected);
socket.on("error", disconnected);
socket.on("reconnect_error", disconnected);
socket.on("reconnect_failed", disconnected);
socket.on("reconnect_attempt", disconnected);
socket.on("reconnecting", disconnected);

// Listen for messages
socket.on("message", function(message) {
  function showCount(record) {
    var count = message.count,
    formattedCount = (new Number(count)).toLocaleString()

    $("#count").text(formattedCount)
    $("#hostname").text(message.hostname)
  }

  if (message.count < 0) {
    // Negative count means the backend counting service cannot be discovered
    disconnectedFromBackendService()
  } else {
    didConnect()
  }
  showCount(message);
});

function didConnect() {
  $("#connection-status").addClass("connected");
  $("#connection-status").text("Connected");
}

socket.on("connect", function() {
  didConnect()

  // Broadcast a message
  function broadcastMessage() {
    socket.emit("send", {"message":"get count"}, function(result) {
      // Silent success, reload again
      setTimeout(broadcastMessage, 200) // In milliseconds
    });
  }
  broadcastMessage();
});
