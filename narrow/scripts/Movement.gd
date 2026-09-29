extends CharacterBody2D


const SPEED = 200.0
const JUMP_VELOCITY = -400.0
const SEND_INTERVAL = 0.05   # 20 updates per second
const REMOTE_SMOOTHING = 0.25

# Set by Game.gd. Only the local player reads keyboard input.
var is_local := false

# Where a remote player should be (updated from the network)
var target_position := Vector2.ZERO

var send_timer := 0.0
var last_sent_position := Vector2.INF


func _physics_process(delta: float) -> void:
	if not is_local:
		# Smoothly follow the last position received from the server
		position = position.lerp(target_position, REMOTE_SMOOTHING)
		return

	if not is_on_floor():
		velocity += get_gravity() * delta

	if Input.is_action_just_pressed("ui_accept") and is_on_floor():
		velocity.y = JUMP_VELOCITY

	var direction := Input.get_axis("ui_left", "ui_right")
	if direction:
		velocity.x = direction * SPEED
	else:
		velocity.x = move_toward(velocity.x, 0, SPEED)

	move_and_slide()

	# Send our position to the server a few times per second
	send_timer += delta
	if send_timer >= SEND_INTERVAL:
		send_timer = 0.0
		if position != last_sent_position:
			last_sent_position = position
			Network.send_position(position)