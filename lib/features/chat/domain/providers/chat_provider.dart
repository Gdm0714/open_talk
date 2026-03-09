import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/network/api_client.dart';
import '../../../../shared/services/websocket_service.dart';
import '../../data/models/chat_room_model.dart';
import '../../data/models/message_model.dart';
import '../../data/repositories/chat_repository.dart';

final chatRepositoryProvider = Provider<ChatRepository>((ref) {
  return ChatRepository(dio: ref.read(apiClientProvider));
});

final chatListProvider =
    AsyncNotifierProvider<ChatListNotifier, List<ChatRoomModel>>(
  ChatListNotifier.new,
);

class ChatListNotifier extends AsyncNotifier<List<ChatRoomModel>> {
  @override
  Future<List<ChatRoomModel>> build() async {
    return ref.read(chatRepositoryProvider).getChats();
  }

  Future<void> refresh() async {
    state = const AsyncValue.loading();
    state = await AsyncValue.guard(
      () => ref.read(chatRepositoryProvider).getChats(),
    );
  }

  void addChat(ChatRoomModel chat) {
    final current = state.valueOrNull ?? [];
    state = AsyncValue.data([chat, ...current]);
  }
}

final chatMessagesProvider = AutoDisposeAsyncNotifierProviderFamily<
    ChatMessagesNotifier, List<MessageModel>, String>(
  ChatMessagesNotifier.new,
);

class ChatMessagesNotifier
    extends AutoDisposeFamilyAsyncNotifier<List<MessageModel>, String> {
  @override
  Future<List<MessageModel>> build(String arg) async {
    return ref.read(chatRepositoryProvider).getChatMessages(arg);
  }

  Future<void> sendMessage(String content) async {
    final repo = ref.read(chatRepositoryProvider);
    final message = await repo.sendMessage(arg, content);
    final current = state.valueOrNull ?? [];
    state = AsyncValue.data([message, ...current]);
  }

  Future<void> loadMore(int page) async {
    final repo = ref.read(chatRepositoryProvider);
    final messages = await repo.getChatMessages(arg, page: page);
    final current = state.valueOrNull ?? [];
    state = AsyncValue.data([...current, ...messages]);
  }

  void addIncomingMessage(WebSocketMessage wsMsg) {
    final currentState = state.valueOrNull;
    if (currentState == null) return;

    final message = MessageModel(
      id: wsMsg.data?['id'] as String? ?? '',
      chatRoomId: wsMsg.roomId ?? '',
      senderId: wsMsg.senderId ?? '',
      senderNickname: wsMsg.senderName ?? '',
      content: wsMsg.content ?? '',
      messageType: MessageType.text,
      createdAt: DateTime.now(),
    );

    state = AsyncData([message, ...currentState]);
  }
}
