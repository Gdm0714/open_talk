import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../shared/widgets/avatar_widget.dart';
import '../../../../shared/widgets/loading_widget.dart';
import '../../domain/providers/home_provider.dart';

class AddFriendScreen extends ConsumerStatefulWidget {
  const AddFriendScreen({super.key});

  @override
  ConsumerState<AddFriendScreen> createState() => _AddFriendScreenState();
}

class _AddFriendScreenState extends ConsumerState<AddFriendScreen> {
  final _searchController = TextEditingController();
  String _query = '';
  final Set<String> _sentRequests = {};

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  Future<void> _sendRequest(String userId) async {
    try {
      await ref.read(friendListProvider.notifier).sendRequest(userId);
      setState(() => _sentRequests.add(userId));
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('친구 요청을 보냈습니다')),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('친구 요청 실패: $e')),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final searchAsync = ref.watch(searchUsersProvider(_query));

    return Scaffold(
      appBar: AppBar(title: const Text('친구 추가')),
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(16),
            child: TextField(
              controller: _searchController,
              autofocus: true,
              onChanged: (value) => setState(() => _query = value.trim()),
              decoration: InputDecoration(
                hintText: '이메일 또는 닉네임으로 검색',
                prefixIcon:
                    const Icon(Icons.search, color: AppColors.textHint),
                suffixIcon: _query.isNotEmpty
                    ? IconButton(
                        icon: const Icon(Icons.clear),
                        onPressed: () {
                          _searchController.clear();
                          setState(() => _query = '');
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
          Expanded(
            child: _query.isEmpty
                ? Center(
                    child: Text(
                      '이메일 또는 닉네임을 입력하세요',
                      style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                            color: AppColors.textSecondary,
                          ),
                    ),
                  )
                : searchAsync.when(
                    data: (users) {
                      if (users.isEmpty) {
                        return Center(
                          child: Text(
                            '검색 결과가 없습니다',
                            style: Theme.of(context)
                                .textTheme
                                .bodyMedium
                                ?.copyWith(color: AppColors.textSecondary),
                          ),
                        );
                      }
                      return ListView.builder(
                        itemCount: users.length,
                        itemBuilder: (context, index) {
                          final user = users[index];
                          final alreadySent = _sentRequests.contains(user.id);
                          return ListTile(
                            leading: AvatarWidget(
                              name: user.nickname,
                              imageUrl: user.avatarUrl,
                              size: 44,
                            ),
                            title: Text(user.nickname),
                            subtitle: Text(
                              user.email,
                              style: const TextStyle(
                                  color: AppColors.textSecondary),
                            ),
                            trailing: alreadySent
                                ? const Chip(label: Text('요청됨'))
                                : TextButton(
                                    onPressed: () => _sendRequest(user.id),
                                    child: const Text('친구 추가'),
                                  ),
                          );
                        },
                      );
                    },
                    loading: () => const LoadingWidget(),
                    error: (e, _) => Center(
                      child: Text('검색 오류: $e'),
                    ),
                  ),
          ),
        ],
      ),
    );
  }
}
