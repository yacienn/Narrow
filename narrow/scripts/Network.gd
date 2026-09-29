#the communication layer between your Godot game and your Go
extends Node

var socket := WebSocketPeer.new()

const SERVER_URL = "ws://192.168.100.30:8080/game"
var connected := false


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


func send_message(message: String):
	if socket.get_ready_state() != WebSocketPeer.STATE_OPEN:
		print("Cannot send message: not connected")
		return

	socket.send_text(message)
