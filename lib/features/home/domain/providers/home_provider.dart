import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/network/api_client.dart';
import '../../../auth/data/models/user_model.dart';
import '../../data/models/friend_model.dart';
import '../../data/repositories/friend_repository.dart';

final friendRepositoryProvider = Provider<FriendRepository>((ref) {
  return FriendRepository(dio: ref.read(apiClientProvider));
});

final friendListProvider =
    AsyncNotifierProvider<FriendListNotifier, List<FriendModel>>(
  FriendListNotifier.new,
);

class FriendListNotifier extends AsyncNotifier<List<FriendModel>> {
  @override
  Future<List<FriendModel>> build() async {
    return ref.read(friendRepositoryProvider).getFriends();
  }

  Future<void> refresh() async {
    state = const AsyncValue.loading();
    state = await AsyncValue.guard(
      () => ref.read(friendRepositoryProvider).getFriends(),
    );
  }

  Future<void> sendRequest(String userId) async {
    await ref.read(friendRepositoryProvider).sendRequest(userId);
  }

  Future<void> acceptRequest(String id) async {
    await ref.read(friendRepositoryProvider).acceptRequest(id);
    await refresh();
  }

  Future<void> rejectRequest(String id) async {
    await ref.read(friendRepositoryProvider).rejectRequest(id);
    await refresh();
  }
}

final searchUsersProvider =
    FutureProvider.autoDispose.family<List<UserModel>, String>(
  (ref, query) async {
    if (query.trim().isEmpty) return [];
    return ref.read(friendRepositoryProvider).searchUsers(query);
  },
);
