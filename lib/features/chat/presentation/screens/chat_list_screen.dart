import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../shared/widgets/avatar_widget.dart';
import '../../../../shared/widgets/loading_widget.dart';
import '../../../auth/domain/providers/auth_provider.dart';
import '../../../home/data/models/friend_model.dart';
import '../../../home/domain/providers/home_provider.dart';
import '../../domain/providers/chat_provider.dart';
import '../widgets/chat_room_tile.dart';

class ChatListScreen extends ConsumerWidget {
  const ChatListScreen({super.key});

  Future<void> _showNewChatDialog(BuildContext context, WidgetRef ref) async {
    final friendListAsync = ref.read(friendListProvider);
    final friends = friendListAsync.valueOrNull
            ?.where((f) => f.status == FriendStatus.accepted)
            .toList() ??
        [];

    if (friends.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('채팅을 시작할 친구가 없습니다'),
          behavior: SnackBarBehavior.floating,
        ),
      );
      return;
    }

    await showModalBottomSheet<void>(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (ctx) => _FriendPickerSheet(friends: friends),
    );
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final chatListAsync = ref.watch(chatListProvider);
    final currentUser = ref.watch(currentUserProvider);
    final currentUserId = currentUser?.id ?? '';

    return Scaffold(
      appBar: AppBar(
        title: const Text('채팅'),
        actions: [
          IconButton(
            icon: const Icon(Icons.search),
            onPressed: () {},
          ),
        ],
      ),
      body: chatListAsync.when(
        data: (chats) {
          if (chats.isEmpty) {
            return Center(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Icon(
                    Icons.chat_bubble_outline,
                    size: 64,
                    color: AppColors.textHint.withValues(alpha: 0.5),
                  ),
                  const SizedBox(height: 16),
                  Text(
                    '아직 대화가 없습니다',
                    style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                          color: AppColors.textSecondary,
                        ),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    '친구에게 메시지를 보내보세요!',
                    style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                          color: AppColors.textHint,
                        ),
                  ),
                ],
              ),
            );
          }

          return RefreshIndicator(
            onRefresh: () => ref.read(chatListProvider.notifier).refresh(),
            child: ListView.separated(
              itemCount: chats.length,
              separatorBuilder: (context, index) => const Divider(indent: 82),
              itemBuilder: (context, index) {
                final chat = chats[index];
                return ChatRoomTile(
                  chatRoom: chat,
                  currentUserId: currentUserId,
                  onTap: () => context.push('/chat/${chat.id}'),
                );
              },
            ),
          );
        },
        loading: () => const LoadingWidget(),
        error: (error, _) => Center(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              const Icon(
                Icons.error_outline,
                size: 48,
                color: AppColors.error,
              ),
              const SizedBox(height: 16),
              Text(
                '채팅 목록을 불러올 수 없습니다',
                style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                      color: AppColors.textSecondary,
                    ),
              ),
              const SizedBox(height: 16),
              TextButton.icon(
                onPressed: () =>
                    ref.read(chatListProvider.notifier).refresh(),
                icon: const Icon(Icons.refresh),
                label: const Text('다시 시도'),
              ),
            ],
          ),
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showNewChatDialog(context, ref),
        child: const Icon(Icons.edit),
      ),
    );
  }
}

class _FriendPickerSheet extends ConsumerWidget {
  final List<FriendModel> friends;

  const _FriendPickerSheet({required this.friends});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return DraggableScrollableSheet(
      initialChildSize: 0.6,
      minChildSize: 0.4,
      maxChildSize: 0.9,
      expand: false,
      builder: (ctx, scrollController) {
        return Column(
          children: [
            const SizedBox(height: 12),
            Container(
              width: 40,
              height: 4,
              decoration: BoxDecoration(
                color: AppColors.divider,
                borderRadius: BorderRadius.circular(2),
              ),
            ),
            const SizedBox(height: 16),
            Text(
              '대화 상대 선택',
              style: Theme.of(context).textTheme.titleMedium,
            ),
            const SizedBox(height: 8),
            const Divider(),
            Expanded(
              child: ListView.builder(
                controller: scrollController,
                itemCount: friends.length,
                itemBuilder: (context, index) {
                  final friend = friends[index];
                  return ListTile(
                    leading: AvatarWidget(
                      name: friend.friendNickname,
                      imageUrl: friend.friendAvatarUrl,
                      size: 44,
                    ),
                    title: Text(friend.friendNickname),
                    subtitle: friend.friendStatusMessage != null &&
                            friend.friendStatusMessage!.isNotEmpty
                        ? Text(
                            friend.friendStatusMessage!,
                            style: const TextStyle(
                                color: AppColors.textSecondary),
                          )
                        : null,
                    onTap: () async {
                      Navigator.of(context).pop();
                      try {
                        final chatRepo = ref.read(chatRepositoryProvider);
                        final chatRoom =
                            await chatRepo.createDirectChat(friend.friendId);
                        ref
                            .read(chatListProvider.notifier)
                            .addChat(chatRoom);
                        if (context.mounted) {
                          context.push('/chat/${chatRoom.id}');
                        }
                      } catch (e) {
                        if (context.mounted) {
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(
                              content: Text('채팅방을 만들 수 없습니다'),
                              backgroundColor: AppColors.error,
                              behavior: SnackBarBehavior.floating,
                            ),
                          );
                        }
                      }
                    },
                  );
                },
              ),
            ),
          ],
        );
      },
    );
  }
}
