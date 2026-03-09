import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../shared/widgets/avatar_widget.dart';
import '../../../auth/domain/providers/auth_provider.dart';

class ProfileScreen extends ConsumerWidget {
  const ProfileScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final currentUser = ref.watch(currentUserProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('설정'),
      ),
      body: ListView(
        children: [
          const SizedBox(height: 16),

          // Profile card
          Container(
            margin: const EdgeInsets.symmetric(horizontal: 16),
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              color: AppColors.surface,
              borderRadius: BorderRadius.circular(16),
            ),
            child: Column(
              children: [
                AvatarWidget(
                  name: currentUser?.nickname ?? '?',
                  imageUrl: currentUser?.avatarUrl,
                  size: 80,
                ),
                const SizedBox(height: 16),
                Text(
                  currentUser?.nickname ?? '',
                  style: Theme.of(context).textTheme.titleLarge,
                ),
                const SizedBox(height: 4),
                Text(
                  currentUser?.email ?? '',
                  style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                        color: AppColors.textSecondary,
                      ),
                ),
                if (currentUser?.statusMessage != null &&
                    currentUser!.statusMessage!.isNotEmpty) ...[
                  const SizedBox(height: 8),
                  Text(
                    currentUser.statusMessage!,
                    style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                          color: AppColors.textSecondary,
                        ),
                    textAlign: TextAlign.center,
                  ),
                ],
                const SizedBox(height: 16),
                SizedBox(
                  width: double.infinity,
                  child: OutlinedButton.icon(
                    onPressed: () => context.push('/profile/edit'),
                    icon: const Icon(Icons.edit, size: 18),
                    label: const Text('프로필 편집'),
                  ),
                ),
              ],
            ),
          ),

          const SizedBox(height: 24),

          // Settings section
          _buildSectionHeader(context, '일반'),
          _buildSettingsTile(
            context,
            icon: Icons.notifications_outlined,
            title: '알림',
            onTap: () {},
          ),
          _buildSettingsTile(
            context,
            icon: Icons.lock_outline,
            title: '개인정보 보호',
            onTap: () {},
          ),
          _buildSettingsTile(
            context,
            icon: Icons.key_outlined,
            title: '비밀번호 변경',
            onTap: () => context.push('/password/change'),
          ),
          _buildSettingsTile(
            context,
            icon: Icons.palette_outlined,
            title: '테마',
            onTap: () {},
          ),

          const SizedBox(height: 16),
          _buildSectionHeader(context, '기타'),
          _buildSettingsTile(
            context,
            icon: Icons.help_outline,
            title: '도움말',
            onTap: () {},
          ),
          _buildSettingsTile(
            context,
            icon: Icons.info_outline,
            title: '앱 정보',
            onTap: () {},
          ),

          const SizedBox(height: 24),

          // Logout
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: TextButton.icon(
              onPressed: () async {
                final confirmed = await showDialog<bool>(
                  context: context,
                  builder: (context) => AlertDialog(
                    title: const Text('로그아웃'),
                    content: const Text('로그아웃 하시겠습니까?'),
                    actions: [
                      TextButton(
                        onPressed: () => Navigator.pop(context, false),
                        child: const Text('취소'),
                      ),
                      TextButton(
                        onPressed: () => Navigator.pop(context, true),
                        child: const Text(
                          '로그아웃',
                          style: TextStyle(color: AppColors.error),
                        ),
                      ),
                    ],
                  ),
                );

                if (confirmed == true) {
                  await ref.read(authStateProvider.notifier).logout();
                  if (context.mounted) {
                    context.go('/login');
                  }
                }
              },
              icon: const Icon(Icons.logout, color: AppColors.error),
              label: const Text(
                '로그아웃',
                style: TextStyle(color: AppColors.error),
              ),
            ),
          ),

          const SizedBox(height: 32),
        ],
      ),
    );
  }

  Widget _buildSectionHeader(BuildContext context, String title) {
    return Padding(
      padding: const EdgeInsets.only(left: 16, bottom: 4),
      child: Text(
        title,
        style: Theme.of(context).textTheme.bodySmall?.copyWith(
              color: AppColors.textHint,
              fontWeight: FontWeight.w600,
            ),
      ),
    );
  }

  Widget _buildSettingsTile(
    BuildContext context, {
    required IconData icon,
    required String title,
    VoidCallback? onTap,
  }) {
    return ListTile(
      leading: Icon(icon, color: AppColors.textSecondary),
      title: Text(title),
      trailing: const Icon(
        Icons.chevron_right,
        color: AppColors.textHint,
      ),
      onTap: onTap,
    );
  }
}
