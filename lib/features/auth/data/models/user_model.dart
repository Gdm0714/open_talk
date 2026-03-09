class UserModel {
  final String id;
  final String email;
  final String nickname;
  final String? avatarUrl;
  final String? statusMessage;

  const UserModel({
    required this.id,
    required this.email,
    required this.nickname,
    this.avatarUrl,
    this.statusMessage,
  });

  factory UserModel.fromJson(Map<String, dynamic> json) {
    return UserModel(
      id: json['id'] as String,
      email: json['email'] as String,
      nickname: json['nickname'] as String,
      avatarUrl: json['avatar_url'] as String?,
      statusMessage: json['status_message'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'email': email,
      'nickname': nickname,
      'avatar_url': avatarUrl,
      'status_message': statusMessage,
    };
  }

  UserModel copyWith({
    String? id,
    String? email,
    String? nickname,
    String? avatarUrl,
    String? statusMessage,
  }) {
    return UserModel(
      id: id ?? this.id,
      email: email ?? this.email,
      nickname: nickname ?? this.nickname,
      avatarUrl: avatarUrl ?? this.avatarUrl,
      statusMessage: statusMessage ?? this.statusMessage,
    );
  }
}
