extends Node2D

@export var player_scene: PackedScene

@onready var players = $Players


func _ready():
	Network.connected_to_server.connect(_on_connected)
	Network.player_joined.connect(_on_player_joined)
	Network.player_left.connect(_on_player_left)

	# In case messages arrived before this scene finished loading
	if Network.player_id != -1:
		spawn_player(Network.player_id, true)

	for id in Network.known_players:
		spawn_player(id, false)


func _on_connected():
	spawn_player(Network.player_id, true)


func _on_player_joined(player_id: int):
	print("Spawning player: ", player_id)
	spawn_player(player_id, false)


func _on_player_left(player_id: int):
	print("Removing player: ", player_id)

	var player = players.get_node_or_null("Player_" + str(player_id))

	if player:
		player.queue_free()


func spawn_player(player_id: int, is_local: bool):
	# Don't spawn the same player twice
	if players.has_node("Player_" + str(player_id)):
		return

	var player = player_scene.instantiate()

	player.name = "Player_" + str(player_id)
	player.is_local = is_local

	players.add_child(player)

	player.position = Vector2(player_id * 40.0, 0)

	print("Spawned Player ", player_id, " at ", player.position, " local=", is_local)