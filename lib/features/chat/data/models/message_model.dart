enum MessageType { text, image, system }

class MessageModel {
  final String id;
  final String chatRoomId;
  final String senderId;
  final String senderNickname;
  final String content;
  final MessageType messageType;
  final DateTime createdAt;

  const MessageModel({
    required this.id,
    required this.chatRoomId,
    required this.senderId,
    required this.senderNickname,
    required this.content,
    this.messageType = MessageType.text,
    required this.createdAt,
  });

  factory MessageModel.fromJson(Map<String, dynamic> json) {
    return MessageModel(
      id: json['id'] as String,
      chatRoomId: json['chat_room_id'] as String,
      senderId: json['sender_id'] as String?
          ?? (json['sender'] != null
              ? (json['sender'] as Map<String, dynamic>)['id'] as String?
              : null)
          ?? '',
      senderNickname: json['sender_nickname'] as String?
          ?? (json['sender'] != null
              ? (json['sender'] as Map<String, dynamic>)['nickname'] as String?
              : null)
          ?? '',
      content: json['content'] as String,
      messageType: _parseMessageType(json['message_type'] as String?),
      createdAt: DateTime.parse(json['created_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'chat_room_id': chatRoomId,
      'sender_id': senderId,
      'sender_nickname': senderNickname,
      'content': content,
      'message_type': messageType.name,
      'created_at': createdAt.toIso8601String(),
    };
  }

  static MessageType _parseMessageType(String? type) {
    switch (type) {
      case 'image':
        return MessageType.image;
      case 'system':
        return MessageType.system;
      default:
        return MessageType.text;
    }
  }

  bool isMine(String currentUserId) => senderId == currentUserId;
}
