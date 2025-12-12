import 'user_model.dart';

class AuthResponse {
  final bool success;
  final String? message;
  final AuthData? data;
  final ErrorData? error;

  AuthResponse({
    required this.success,
    this.message,
    this.data,
    this.error,
  });

  factory AuthResponse.fromJson(Map<String, dynamic> json) {
    return AuthResponse(
      success: json['success'] as bool? ?? false,
      message: json['message'] as String?,
      data: json['data'] != null ? AuthData.fromJson(json['data'] as Map<String, dynamic>) : null,
      error: json['error'] != null ? ErrorData.fromJson(json['error'] as Map<String, dynamic>) : null,
    );
  }
}

class AuthData {
  final String token;
  final UserModel user;

  AuthData({
    required this.token,
    required this.user,
  });

  factory AuthData.fromJson(Map<String, dynamic> json) {
    return AuthData(
      token: json['token'] as String,
      user: UserModel.fromJson(json['user'] as Map<String, dynamic>),
    );
  }
}

class ErrorData {
  final String code;
  final String message;

  ErrorData({
    required this.code,
    required this.message,
  });

  factory ErrorData.fromJson(Map<String, dynamic> json) {
    return ErrorData(
      code: json['code'] as String? ?? 'UNKNOWN_ERROR',
      message: json['message'] as String? ?? 'An unknown error occurred',
    );
  }
}
