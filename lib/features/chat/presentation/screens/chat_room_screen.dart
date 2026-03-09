import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../shared/services/websocket_service.dart';
import '../../../../shared/widgets/loading_widget.dart';
import '../../../auth/domain/providers/auth_provider.dart';
import '../../domain/providers/chat_provider.dart';
import '../widgets/message_bubble.dart';

class ChatRoomScreen extends ConsumerStatefulWidget {
  final String roomId;

  const ChatRoomScreen({super.key, required this.roomId});

  @override
  ConsumerState<ChatRoomScreen> createState() => _ChatRoomScreenState();
}

class _ChatRoomScreenState extends ConsumerState<ChatRoomScreen> {
  final _messageController = TextEditingController();
  final _scrollController = ScrollController();
  bool _isSending = false;
  StreamSubscription<WebSocketMessage>? _wsSubscription;
  String? _typingUser;
  Timer? _typingTimer;
  Timer? _typingDebounce;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _initWebSocket();
      _markAsRead();
    });
  }

  void _initWebSocket() {
    final wsService = ref.read(webSocketServiceProvider);
    wsService.joinRoom(widget.roomId);

    _wsSubscription = wsService.messageStream.listen((msg) {
      if (!mounted) return;
      final currentUserId = ref.read(currentUserProvider)?.id ?? '';

      if (msg.type == 'message' && msg.roomId == widget.roomId) {
        // Only add if not from current user (our own messages are added optimistically)
        if (msg.senderId != currentUserId) {
          ref
              .read(chatMessagesProvider(widget.roomId).notifier)
              .addIncomingMessage(msg);
        }
      }

      if (msg.type == 'typing' &&
          msg.roomId == widget.roomId &&
          msg.senderId != currentUserId) {
        setState(() => _typingUser = msg.senderName);
        _typingTimer?.cancel();
        _typingTimer = Timer(const Duration(seconds: 3), () {
          if (mounted) setState(() => _typingUser = null);
        });
      }
    });
  }

  Future<void> _markAsRead() async {
    try {
      final repo = ref.read(chatRepositoryProvider);
      await repo.markAsRead(widget.roomId);
    } catch (_) {
      // Best effort
    }
  }

  @override
  void dispose() {
    _wsSubscription?.cancel();
    _typingTimer?.cancel();
    _typingDebounce?.cancel();
    final wsService = ref.read(webSocketServiceProvider);
    wsService.leaveRoom(widget.roomId);
    _messageController.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  Future<void> _sendMessage() async {
    final content = _messageController.text.trim();
    if (content.isEmpty || _isSending) return;

    setState(() => _isSending = true);
    _messageController.clear();

    try {
      await ref
          .read(chatMessagesProvider(widget.roomId).notifier)
          .sendMessage(content);
    } catch (_) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('메시지 전송에 실패했습니다'),
            backgroundColor: AppColors.error,
            behavior: SnackBarBehavior.floating,
          ),
        );
      }
    } finally {
      if (mounted) setState(() => _isSending = false);
    }
  }

  void _onTextChanged(String text) {
    _typingDebounce?.cancel();
    _typingDebounce = Timer(const Duration(milliseconds: 500), () {
      final wsService = ref.read(webSocketServiceProvider);
      wsService.sendTyping(widget.roomId);
    });
  }

  @override
  Widget build(BuildContext context) {
    final messagesAsync = ref.watch(chatMessagesProvider(widget.roomId));
    final currentUser = ref.watch(currentUserProvider);
    final currentUserId = currentUser?.id ?? '';

    // Get chat room info from chat list
    final chatList = ref.watch(chatListProvider).valueOrNull ?? [];
    final chatRoom = chatList.where((c) => c.id == widget.roomId).firstOrNull;
    final roomName = chatRoom?.displayName(currentUserId) ?? '채팅';
    final memberCount = chatRoom?.members.length ?? 0;

    return Scaffold(
      appBar: AppBar(
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              roomName,
              style: const TextStyle(fontSize: 17, fontWeight: FontWeight.w600),
            ),
            if (memberCount > 0)
              Text(
                '$memberCount명',
                style: TextStyle(
                  fontSize: 12,
                  color: AppColors.textSecondary,
                  fontWeight: FontWeight.w400,
                ),
              ),
          ],
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.menu),
            onPressed: () {},
          ),
        ],
      ),
      body: Column(
        children: [
          // Messages
          Expanded(
            child: messagesAsync.when(
              data: (messages) {
                if (messages.isEmpty) {
                  return Center(
                    child: Text(
                      '첫 번째 메시지를 보내보세요!',
                      style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                            color: AppColors.textHint,
                          ),
                    ),
                  );
                }

                return ListView.builder(
                  controller: _scrollController,
                  reverse: true,
                  padding: const EdgeInsets.symmetric(vertical: 8),
                  itemCount: messages.length,
                  itemBuilder: (context, index) {
                    final message = messages[index];
                    final isMine = message.isMine(currentUserId);
                    final showSender = !isMine &&
                        (index == messages.length - 1 ||
                            messages[index + 1].senderId !=
                                message.senderId);

                    return MessageBubble(
                      message: message,
                      isMine: isMine,
                      showSenderName: showSender,
                    );
                  },
                );
              },
              loading: () => const LoadingWidget(),
              error: (error, _) => Center(
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    const Icon(Icons.error_outline, color: AppColors.error),
                    const SizedBox(height: 8),
                    const Text('메시지를 불러올 수 없습니다'),
                    TextButton(
                      onPressed: () => ref.invalidate(
                        chatMessagesProvider(widget.roomId),
                      ),
                      child: const Text('다시 시도'),
                    ),
                  ],
                ),
              ),
            ),
          ),

          // Typing indicator
          if (_typingUser != null)
            Padding(
              padding: const EdgeInsets.only(left: 16, bottom: 4),
              child: Align(
                alignment: Alignment.centerLeft,
                child: Text(
                  '$_typingUser 님이 입력 중...',
                  style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        color: AppColors.textHint,
                        fontStyle: FontStyle.italic,
                      ),
                ),
              ),
            ),

          // Input bar
          _buildInputBar(context),
        ],
      ),
    );
  }

  Widget _buildInputBar(BuildContext context) {
    return Container(
      padding: EdgeInsets.only(
        left: 12,
        right: 8,
        top: 8,
        bottom: MediaQuery.of(context).padding.bottom + 8,
      ),
      decoration: const BoxDecoration(
        color: AppColors.surface,
        border: Border(
          top: BorderSide(color: AppColors.divider, width: 0.5),
        ),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.end,
        children: [
          IconButton(
            icon: const Icon(Icons.add, color: AppColors.textSecondary),
            onPressed: () {},
          ),
          Expanded(
            child: Container(
              constraints: const BoxConstraints(maxHeight: 120),
              child: TextField(
                controller: _messageController,
                maxLines: null,
                textInputAction: TextInputAction.newline,
                onChanged: _onTextChanged,
                decoration: InputDecoration(
                  hintText: '메시지 보내기',
                  hintStyle: const TextStyle(color: AppColors.textHint),
                  filled: true,
                  fillColor: AppColors.surfaceVariant,
                  contentPadding: const EdgeInsets.symmetric(
                    horizontal: 16,
                    vertical: 10,
                  ),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(24),
                    borderSide: BorderSide.none,
                  ),
                  enabledBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(24),
                    borderSide: BorderSide.none,
                  ),
                  focusedBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(24),
                    borderSide: BorderSide.none,
                  ),
                ),
              ),
            ),
          ),
          const SizedBox(width: 4),
          Container(
            decoration: const BoxDecoration(
              color: AppColors.primary,
              shape: BoxShape.circle,
            ),
            child: IconButton(
              icon: _isSending
                  ? const SizedBox(
                      width: 20,
                      height: 20,
                      child: CircularProgressIndicator(
                        strokeWidth: 2,
                        color: Colors.white,
                      ),
                    )
                  : const Icon(Icons.send, color: Colors.white, size: 20),
              onPressed: _sendMessage,
            ),
          ),
        ],
      ),
    );
  }
}
