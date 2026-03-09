class ApiEndpoints {
  ApiEndpoints._();

  // Auth
  static const String login = '/auth/login';
  static const String register = '/auth/register';
  static const String refreshToken = '/auth/refresh';
  static const String logout = '/auth/logout';

  // Users
  static const String me = '/users/me';
  static const String users = '/users';
  static String userById(String id) => '/users/$id';
  static String searchUsers(String query) => '/users/search?q=$query';

  // Friends
  static const String friends = '/friends';
  static String friendById(String id) => '/friends/$id';
  static const String friendRequests = '/friends/requests';
  static String acceptFriend(String id) => '/friends/$id/accept';
  static String rejectFriend(String id) => '/friends/$id/reject';

  // Chats
  static const String chats = '/chats';
  static String chatById(String id) => '/chats/$id';
  static String chatMessages(String roomId) => '/chats/$roomId/messages';
}
