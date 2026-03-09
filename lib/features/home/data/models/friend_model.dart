enum FriendStatus { accepted, pending, rejected }

class FriendModel {
  final String id;
  final String friendId;
  final String friendNickname;
  final String? friendAvatarUrl;
  final String? friendStatusMessage;
  final FriendStatus status;

  const FriendModel({
    required this.id,
    required this.friendId,
    required this.friendNickname,
    this.friendAvatarUrl,
    this.friendStatusMessage,
    this.status = FriendStatus.accepted,
  });

  factory FriendModel.fromJson(Map<String, dynamic> json) {
    return FriendModel(
      id: json['id'] as String,
      friendId: json['friend_id'] as String,
      friendNickname: json['friend_nickname'] as String,
      friendAvatarUrl: json['friend_avatar_url'] as String?,
      friendStatusMessage: json['friend_status_message'] as String?,
      status: _parseStatus(json['status'] as String?),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'friend_id': friendId,
      'friend_nickname': friendNickname,
      'friend_avatar_url': friendAvatarUrl,
      'friend_status_message': friendStatusMessage,
      'status': status.name,
    };
  }

  static FriendStatus _parseStatus(String? status) {
    switch (status) {
      case 'pending':
        return FriendStatus.pending;
      case 'rejected':
        return FriendStatus.rejected;
      default:
        return FriendStatus.accepted;
    }
  }
}
