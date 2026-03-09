import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../shared/widgets/loading_widget.dart';
import '../../../chat/domain/providers/chat_provider.dart';
import '../../data/models/friend_model.dart';
import '../../domain/providers/home_provider.dart';
import '../widgets/friend_tile.dart';

class FriendListScreen extends ConsumerStatefulWidget {
  const FriendListScreen({super.key});

  @override
  ConsumerState<FriendListScreen> createState() => _FriendListScreenState();
}

class _FriendListScreenState extends ConsumerState<FriendListScreen> {
  final _searchController = TextEditingController();
  String _searchQuery = '';

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  Future<void> _startDirectChat(FriendModel friend) async {
    try {
      final chatRepo = ref.read(chatRepositoryProvider);
      final chatRoom = await chatRepo.createDirectChat(friend.friendId);
      ref.read(chatListProvider.notifier).addChat(chatRoom);
      if (mounted) {
        context.push('/chat/${chatRoom.id}');
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('채팅방을 만들 수 없습니다'),
            backgroundColor: AppColors.error,
            behavior: SnackBarBehavior.floating,
          ),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final friendListAsync = ref.watch(friendListProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('친구'),
        actions: [
          IconButton(
            icon: const Icon(Icons.person_add_outlined),
            onPressed: () => context.push('/friends/add'),
          ),
          IconButton(
            icon: const Icon(Icons.inbox_outlined),
            onPressed: () => context.push('/friends/requests'),
          ),
        ],
      ),
      body: Column(
        children: [
          // Search bar
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            child: TextField(
              controller: _searchController,
              onChanged: (value) => setState(() => _searchQuery = value),
              decoration: InputDecoration(
                hintText: '친구 검색',
                prefixIcon: const Icon(Icons.search, color: AppColors.textHint),
                suffixIcon: _searchQuery.isNotEmpty
                    ? IconButton(
                        icon: const Icon(Icons.clear),
                        onPressed: () {
                          _searchController.clear();
                          setState(() => _searchQuery = '');
                        },
                      )
                    : null,
                filled: true,
                fillColor: AppColors.surfaceVariant,
                contentPadding: const EdgeInsets.symmetric(vertical: 10),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(24),
                  borderSide: BorderSide.none,
                ),
              ),
            ),
          ),

          // Friend list
          Expanded(
            child: friendListAsync.when(
              data: (friends) {
                final filtered = _searchQuery.isEmpty
                    ? friends
                        .where((f) => f.status == FriendStatus.accepted)
                        .toList()
                    : friends
                        .where((f) =>
                            f.status == FriendStatus.accepted &&
                            f.friendNickname
                                .toLowerCase()
                                .contains(_searchQuery.toLowerCase()))
                        .toList();

                if (filtered.isEmpty) {
                  return Center(
                    child: Column(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Icon(
                          Icons.people_outline,
                          size: 64,
                          color: AppColors.textHint.withValues(alpha: 0.5),
                        ),
                        const SizedBox(height: 16),
                        Text(
                          _searchQuery.isEmpty
                              ? '아직 친구가 없습니다'
                              : '검색 결과가 없습니다',
                          style:
                              Theme.of(context).textTheme.bodyLarge?.copyWith(
                                    color: AppColors.textSecondary,
                                  ),
                        ),
                      ],
                    ),
                  );
                }

                return RefreshIndicator(
                  onRefresh: () =>
                      ref.read(friendListProvider.notifier).refresh(),
                  child: ListView.builder(
                    itemCount: filtered.length + 1,
                    itemBuilder: (context, index) {
                      if (index == 0) {
                        return Padding(
                          padding: const EdgeInsets.only(
                            left: 16,
                            top: 12,
                            bottom: 4,
                          ),
                          child: Text(
                            '친구 ${filtered.length}',
                            style:
                                Theme.of(context).textTheme.bodySmall?.copyWith(
                                      color: AppColors.textHint,
                                    ),
                          ),
                        );
                      }

                      final friend = filtered[index - 1];
                      return FriendTile(
                        friend: friend,
                        onTap: () => _startDirectChat(friend),
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
                    const Icon(Icons.error_outline, color: AppColors.error),
                    const SizedBox(height: 8),
                    const Text('친구 목록을 불러올 수 없습니다'),
                    TextButton(
                      onPressed: () =>
                          ref.read(friendListProvider.notifier).refresh(),
                      child: const Text('다시 시도'),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
