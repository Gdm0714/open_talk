class AppConstants {
  AppConstants._();

  static const String appName = 'OpenTalk';

  // API
  static const String apiBaseUrl = 'http://192.168.0.6:8081/api';
  static const String wsBaseUrl = 'ws://192.168.0.6:8081/ws';

  // Storage Keys
  static const String accessTokenKey = 'access_token';
  static const String refreshTokenKey = 'refresh_token';
  static const String userIdKey = 'user_id';

  // Pagination
  static const int defaultPageSize = 20;
}
