# Communication layer between the Godot game and the Go server

extends Node

var socket := WebSocketPeer.new()

const SERVER_URL = "ws://192.168.100.30:8080/game"

var connected := false
var player_id := -1


# Signals
signal connected_to_server
signal player_joined(player_id)


func _ready():
	connect_to_server()


func connect_to_server():
	var error = socket.connect_to_url(SERVER_URL)

	if error != OK:
		print("Failed to connect: ", error)
		return

	print("Connecting to server...")


func _process(_delta):
	socket.poll()

	var state = socket.get_ready_state()

	if state == WebSocketPeer.STATE_OPEN:

		if not connected:
			connected = true
			print("Connected to Go server!")

		receive_messages()

	elif state == WebSocketPeer.STATE_CLOSED:

		if connected:
			connected = false
			print("Disconnected from Go server")


func receive_messages():

	while socket.get_available_packet_count() > 0:

		var packet = socket.get_packet()
		var message = packet.get_string_from_utf8()

		print("Server: ", message)

		handle_message(message)


func handle_message(message: String):

	var data = JSON.parse_string(message)

	if data == null:
		print("Invalid JSON from server")
		return

	if not data is Dictionary:
		print("Invalid message format")
		return

	var type = data.get("type", "")

	match type:

		"welcome":
			handle_welcome(data)

		"player_joined":
			handle_player_joined(data)

		_:
			print("Unknown message type: ", type)


func handle_welcome(data: Dictionary):

	player_id = data.get("player_id", -1)

	print("My player ID: ", player_id)

	# Tell Game.gd that we are connected
	connected_to_server.emit()


func handle_player_joined(data: Dictionary):

	var joined_player_id = data.get("player_id", -1)

	print("Player joined: ", joined_player_id)

	# Tell Game.gd to spawn this player
	player_joined.emit(joined_player_id)


func send_message(message: String):

	if socket.get_ready_state() != WebSocketPeer.STATE_OPEN:
		print("Cannot send message: not connected")
		return

	socket.send_text(message)