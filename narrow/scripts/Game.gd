extends Node3D

@export var player_scene: PackedScene

@onready var players = $Players


func _ready():
	Network.player_joined.connect(_on_player_joined)
	Network.player_left.connect(_on_player_left)

	if Network.player_id != -1:
		spawn_player(Network.player_id)


func _on_player_joined(player_id: int):
	print("Spawning player: ", player_id)

	spawn_player(player_id)


func _on_player_left(player_id: int):
	print("Removing player: ", player_id)

	var player = players.get_node_or_null("Player_" + str(player_id))

	if player:
		player.queue_free()


func spawn_player(player_id: int):

	var player = player_scene.instantiate()

	player.name = "Player_" + str(player_id)

	players.add_child(player)

	player.position = Vector3(
		player_id * 2.0,
		0,
		0
	)

	print(
		"Spawned Player ",
		player_id,
		" at ",
		player.position
	)