import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/network/api_client.dart';
import '../../../../shared/widgets/app_button.dart';
import '../../../../shared/widgets/app_text_field.dart';
import '../../../auth/domain/providers/auth_provider.dart';

class ProfileEditScreen extends ConsumerStatefulWidget {
  const ProfileEditScreen({super.key});

  @override
  ConsumerState<ProfileEditScreen> createState() => _ProfileEditScreenState();
}

class _ProfileEditScreenState extends ConsumerState<ProfileEditScreen> {
  final _formKey = GlobalKey<FormState>();
  late TextEditingController _nicknameController;
  late TextEditingController _statusMessageController;
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    final user = ref.read(currentUserProvider);
    _nicknameController = TextEditingController(text: user?.nickname ?? '');
    _statusMessageController =
        TextEditingController(text: user?.statusMessage ?? '');
  }

  @override
  void dispose() {
    _nicknameController.dispose();
    _statusMessageController.dispose();
    super.dispose();
  }

  Future<void> _save() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() => _isLoading = true);
    try {
      final dio = ref.read(apiClientProvider);
      await dio.put('/users/me', data: {
        'nickname': _nicknameController.text.trim(),
        'status_message': _statusMessageController.text.trim(),
      });

      if (!mounted) return;
      await ref.read(authStateProvider.notifier).checkAuthStatus();
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('프로필이 수정되었습니다')),
      );
      Navigator.of(context).pop();
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('프로필 수정 실패: $e')),
        );
      }
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('프로필 수정')),
      body: Padding(
        padding: const EdgeInsets.all(24),
        child: Form(
          key: _formKey,
          child: Column(
            children: [
              AppTextField(
                controller: _nicknameController,
                label: '닉네임',
                validator: (v) {
                  if (v == null || v.trim().isEmpty) return '닉네임을 입력하세요';
                  if (v.trim().length < 2) return '닉네임은 2자 이상이어야 합니다';
                  return null;
                },
              ),
              const SizedBox(height: 16),
              AppTextField(
                controller: _statusMessageController,
                label: '상태 메시지',
                validator: (_) => null,
              ),
              const SizedBox(height: 32),
              SizedBox(
                width: double.infinity,
                child: AppButton(
                  text: '저장',
                  onPressed: _isLoading ? null : _save,
                  isLoading: _isLoading,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
