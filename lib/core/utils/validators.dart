class Validators {
  Validators._();

  static String? validateEmail(String? value) {
    if (value == null || value.trim().isEmpty) {
      return '이메일을 입력해주세요';
    }
    final emailRegex = RegExp(r'^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$');
    if (!emailRegex.hasMatch(value.trim())) {
      return '올바른 이메일 형식이 아닙니다';
    }
    return null;
  }

  static String? validatePassword(String? value) {
    if (value == null || value.isEmpty) {
      return '비밀번호를 입력해주세요';
    }
    if (value.length < 8) {
      return '비밀번호는 8자 이상이어야 합니다';
    }
    if (!RegExp(r'[A-Za-z]').hasMatch(value)) {
      return '비밀번호에 영문자를 포함해주세요';
    }
    if (!RegExp(r'[0-9]').hasMatch(value)) {
      return '비밀번호에 숫자를 포함해주세요';
    }
    return null;
  }

  static String? validateNickname(String? value) {
    if (value == null || value.trim().isEmpty) {
      return '닉네임을 입력해주세요';
    }
    if (value.trim().length < 2) {
      return '닉네임은 2자 이상이어야 합니다';
    }
    if (value.trim().length > 20) {
      return '닉네임은 20자 이하여야 합니다';
    }
    return null;
  }

  static String? validateRequired(String? value, String fieldName) {
    if (value == null || value.trim().isEmpty) {
      return '$fieldName을(를) 입력해주세요';
    }
    return null;
  }
}
