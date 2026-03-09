import 'package:dio/dio.dart';

import '../../../../core/network/api_endpoints.dart';
import '../models/chat_room_model.dart';
import '../models/message_model.dart';

class ChatRepository {
  final Dio _dio;

  ChatRepository({required Dio dio}) : _dio = dio;

  Future<List<ChatRoomModel>> getChats() async {
    final response = await _dio.get(ApiEndpoints.chats);
    final data =
        (response.data as Map<String, dynamic>)['data'] as List<dynamic>;
    return data
        .map((item) => ChatRoomModel.fromJson(item as Map<String, dynamic>))
        .toList();
  }

  Future<List<MessageModel>> getChatMessages(
    String roomId, {
    int page = 1,
  }) async {
    final response = await _dio.get(
      ApiEndpoints.chatMessages(roomId),
      queryParameters: {'page': page},
    );
    final data =
        (response.data as Map<String, dynamic>)['data'] as List<dynamic>;
    return data
        .map((item) => MessageModel.fromJson(item as Map<String, dynamic>))
        .toList();
  }

  Future<MessageModel> sendMessage(String roomId, String content) async {
    final response = await _dio.post(
      ApiEndpoints.chatMessages(roomId),
      data: {'content': content, 'message_type': 'text'},
    );
    return MessageModel.fromJson(
      (response.data as Map<String, dynamic>)['data'] as Map<String, dynamic>,
    );
  }

  Future<ChatRoomModel> createDirectChat(String userId) async {
    final response = await _dio.post(
      ApiEndpoints.chats,
      data: {'type': 'direct', 'member_ids': [userId]},
    );
    return ChatRoomModel.fromJson(
      (response.data as Map<String, dynamic>)['data'] as Map<String, dynamic>,
    );
  }

  Future<ChatRoomModel> createGroupChat(
    String name,
    List<String> memberIds,
  ) async {
    final response = await _dio.post(
      ApiEndpoints.chats,
      data: {
        'type': 'group',
        'name': name,
        'member_ids': memberIds,
      },
    );
    return ChatRoomModel.fromJson(
      (response.data as Map<String, dynamic>)['data'] as Map<String, dynamic>,
    );
  }

  Future<void> markAsRead(String roomId) async {
    await _dio.put('/chats/$roomId/read');
  }
}
