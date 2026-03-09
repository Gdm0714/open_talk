import 'package:dio/dio.dart';

import '../../../../core/network/api_endpoints.dart';
import '../../../auth/data/models/user_model.dart';
import '../models/friend_model.dart';

class FriendRepository {
  final Dio _dio;

  FriendRepository({required Dio dio}) : _dio = dio;

  Future<List<FriendModel>> getFriends() async {
    final response = await _dio.get(ApiEndpoints.friends);
    final data =
        (response.data as Map<String, dynamic>)['data'] as List<dynamic>;
    return data
        .map((item) => FriendModel.fromJson(item as Map<String, dynamic>))
        .toList();
  }

  Future<void> sendRequest(String userId) async {
    await _dio.post(
      '${ApiEndpoints.friends}/request',
      data: {'friend_id': userId},
    );
  }

  Future<void> acceptRequest(String id) async {
    await _dio.put(ApiEndpoints.acceptFriend(id));
  }

  Future<void> rejectRequest(String id) async {
    await _dio.put(ApiEndpoints.rejectFriend(id));
  }

  Future<List<UserModel>> searchUsers(String query) async {
    final response = await _dio.get(ApiEndpoints.searchUsers(query));
    final data =
        (response.data as Map<String, dynamic>)['data'] as List<dynamic>;
    return data
        .map((item) => UserModel.fromJson(item as Map<String, dynamic>))
        .toList();
  }
}
