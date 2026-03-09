import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../shared/widgets/avatar_widget.dart';
import '../../../../shared/widgets/loading_widget.dart';
import '../../data/models/friend_model.dart';
import '../../domain/providers/home_provider.dart';

class FriendRequestsScreen extends ConsumerWidget {
  const FriendRequestsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final friendListAsync = ref.watch(friendListProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('친구 요청')),
      body: friendListAsync.when(
        data: (friends) {
          final pending = friends
              .where((f) => f.status == FriendStatus.pending)
              .toList();

          if (pending.isEmpty) {
            return Center(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Icon(
                    Icons.inbox_outlined,
                    size: 64,
                    color: AppColors.textHint.withValues(alpha: 0.5),
                  ),
                  const SizedBox(height: 16),
                  Text(
                    '받은 친구 요청이 없습니다',
                    style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                          color: AppColors.textSecondary,
                        ),
                  ),
                ],
              ),
            );
          }

          return RefreshIndicator(
            onRefresh: () => ref.read(friendListProvider.notifier).refresh(),
            child: ListView.builder(
              itemCount: pending.length,
              itemBuilder: (context, index) {
                final request = pending[index];
                return _FriendRequestTile(request: request);
              },
            ),
          );
        },
        loading: () => const LoadingWidget(),
        error: (error, _) => Center(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              const Icon(Icons.error_outline, color: AppColors.error),
              const SizedBox(height: 8),
              const Text('친구 요청을 불러올 수 없습니다'),
              TextButton(
                onPressed: () => ref.read(friendListProvider.notifier).refresh(),
                child: const Text('다시 시도'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _FriendRequestTile extends ConsumerWidget {
  final FriendModel request;

  const _FriendRequestTile({required this.request});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return ListTile(
      leading: AvatarWidget(
        name: request.friendNickname,
        imageUrl: request.friendAvatarUrl,
        size: 44,
      ),
      title: Text(request.friendNickname),
      subtitle: request.friendStatusMessage != null &&
              request.friendStatusMessage!.isNotEmpty
          ? Text(
              request.friendStatusMessage!,
              style: const TextStyle(color: AppColors.textSecondary),
            )
          : null,
      trailing: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          TextButton(
            onPressed: () async {
              try {
                await ref
                    .read(friendListProvider.notifier)
                    .acceptRequest(request.id);
                if (context.mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('친구 요청을 수락했습니다')),
                  );
                }
              } catch (e) {
                if (context.mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text('수락 실패: $e')),
                  );
                }
              }
            },
            child: const Text('수락'),
          ),
          TextButton(
            onPressed: () async {
              try {
                await ref
                    .read(friendListProvider.notifier)
                    .rejectRequest(request.id);
                if (context.mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('친구 요청을 거절했습니다')),
                  );
                }
              } catch (e) {
                if (context.mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text('거절 실패: $e')),
                  );
                }
              }
            },
            style: TextButton.styleFrom(foregroundColor: AppColors.error),
            child: const Text('거절'),
          ),
        ],
      ),
    );
  }
}
