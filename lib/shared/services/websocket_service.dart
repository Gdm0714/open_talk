import 'dart:async';
import 'dart:convert';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:web_socket_channel/web_socket_channel.dart';
import '../../core/constants/app_constants.dart';

class WebSocketMessage {
  final String type;
  final String? roomId;
  final String? content;
  final String? senderId;
  final String? senderName;
  final Map<String, dynamic>? data;

  WebSocketMessage({
    required this.type,
    this.roomId,
    this.content,
    this.senderId,
    this.senderName,
    this.data,
  });

  factory WebSocketMessage.fromJson(Map<String, dynamic> json) {
    return WebSocketMessage(
      type: json['type'] as String,
      roomId: json['room_id'] as String?,
      content: json['content'] as String?,
      senderId: json['sender_id'] as String?,
      senderName: json['sender_name'] as String?,
      data: json,
    );
  }

  Map<String, dynamic> toJson() => {
        'type': type,
        if (roomId != null) 'room_id': roomId,
        if (content != null) 'content': content,
        if (senderId != null) 'sender_id': senderId,
        if (senderName != null) 'sender_name': senderName,
      };
}

class WebSocketService {
  WebSocketChannel? _channel;
  final _messageController = StreamController<WebSocketMessage>.broadcast();
  final _connectionController = StreamController<bool>.broadcast();
  bool _isConnected = false;
  Timer? _reconnectTimer;
  Timer? _pingTimer;
  String? _token;

  Stream<WebSocketMessage> get messageStream => _messageController.stream;
  Stream<bool> get connectionStream => _connectionController.stream;
  bool get isConnected => _isConnected;

  Future<void> connect(String token) async {
    _token = token;
    await _doConnect();
  }

  Future<void> _doConnect() async {
    if (_token == null) return;

    try {
      final uri = Uri.parse('${AppConstants.wsBaseUrl}?token=$_token');
      _channel = WebSocketChannel.connect(uri);

      _channel!.stream.listen(
        (data) {
          try {
            final json = jsonDecode(data as String) as Map<String, dynamic>;
            final message = WebSocketMessage.fromJson(json);
            _messageController.add(message);
          } catch (e) {
            // Ignore malformed messages
          }
        },
        onDone: () {
          _isConnected = false;
          _connectionController.add(false);
          _scheduleReconnect();
        },
        onError: (error) {
          _isConnected = false;
          _connectionController.add(false);
          _scheduleReconnect();
        },
      );

      _isConnected = true;
      _connectionController.add(true);
      _startPing();
    } catch (e) {
      _isConnected = false;
      _connectionController.add(false);
      _scheduleReconnect();
    }
  }

  void _scheduleReconnect() {
    _reconnectTimer?.cancel();
    _reconnectTimer = Timer(const Duration(seconds: 3), () {
      if (!_isConnected && _token != null) {
        _doConnect();
      }
    });
  }

  void _startPing() {
    _pingTimer?.cancel();
    _pingTimer = Timer.periodic(const Duration(seconds: 30), (_) {
      if (_isConnected) {
        sendRaw({'type': 'ping'});
      }
    });
  }

  void joinRoom(String roomId) {
    sendRaw({'type': 'join', 'room_id': roomId});
  }

  void leaveRoom(String roomId) {
    sendRaw({'type': 'leave', 'room_id': roomId});
  }

  void sendMessage(String roomId, String content) {
    sendRaw({
      'type': 'message',
      'room_id': roomId,
      'content': content,
    });
  }

  void sendTyping(String roomId) {
    sendRaw({
      'type': 'typing',
      'room_id': roomId,
    });
  }

  void sendRaw(Map<String, dynamic> data) {
    if (_isConnected && _channel != null) {
      _channel!.sink.add(jsonEncode(data));
    }
  }

  void disconnect() {
    _reconnectTimer?.cancel();
    _pingTimer?.cancel();
    _channel?.sink.close();
    _isConnected = false;
    _connectionController.add(false);
    _token = null;
  }

  void dispose() {
    disconnect();
    _messageController.close();
    _connectionController.close();
  }
}

final webSocketServiceProvider = Provider<WebSocketService>((ref) {
  final service = WebSocketService();
  ref.onDispose(() => service.dispose());
  return service;
});
