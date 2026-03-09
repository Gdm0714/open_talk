import '../../../auth/data/models/user_model.dart';

enum ChatRoomType { direct, group }

class ChatRoomModel {
  final String id;
  final String? name;
  final ChatRoomType type;
  final String? lastMessage;
  final DateTime? lastMessageAt;
  final List<UserModel> members;
  final int unreadCount;

  const ChatRoomModel({
    required this.id,
    this.name,
    required this.type,
    this.lastMessage,
    this.lastMessageAt,
    this.members = const [],
    this.unreadCount = 0,
  });

  factory ChatRoomModel.fromJson(Map<String, dynamic> json) {
    return ChatRoomModel(
      id: json['id'] as String,
      name: json['name'] as String?,
      type: json['type'] == 'group' ? ChatRoomType.group : ChatRoomType.direct,
      lastMessage: json['last_message'] != null
          ? (json['last_message'] is String
              ? json['last_message'] as String
              : (json['last_message'] as Map<String, dynamic>)['content']
                  as String?)
          : null,
      lastMessageAt: json['last_message'] != null &&
              json['last_message'] is Map
          ? DateTime.tryParse(
              (json['last_message'] as Map<String, dynamic>)['created_at']
                      ?.toString() ??
                  '')
          : (json['last_message_at'] != null
              ? DateTime.tryParse(json['last_message_at'].toString())
              : null),
      members: json['members'] != null
          ? (json['members'] as List<dynamic>)
              .map((m) => UserModel.fromJson(m as Map<String, dynamic>))
              .toList()
          : [],
      unreadCount: json['unread_count'] as int? ?? 0,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'type': type == ChatRoomType.group ? 'group' : 'direct',
      'last_message': lastMessage,
      'last_message_at': lastMessageAt?.toIso8601String(),
      'members': members.map((m) => m.toJson()).toList(),
      'unread_count': unreadCount,
    };
  }

  String displayName(String currentUserId) {
    if (name != null && name!.isNotEmpty) return name!;
    if (type == ChatRoomType.direct) {
      final other = members.where((m) => m.id != currentUserId).firstOrNull;
      return other?.nickname ?? '알 수 없는 사용자';
    }
    return members.map((m) => m.nickname).join(', ');
  }
}
