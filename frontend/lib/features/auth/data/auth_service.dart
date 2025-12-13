import 'package:dio/dio.dart';
import '../../../core/constants/api_constants.dart';
import '../../../core/services/api_service.dart';
import '../../../core/services/storage_service.dart';
import '../domain/models/auth_response.dart';
import '../domain/models/user_model.dart';

class AuthService {
  final Dio _dio = ApiService.dio;

  Future<AuthResponse> login({
    required String email,
    required String password,
  }) async {
    try {
      final response = await _dio.post(
        ApiConstants.login,
        data: {'email': email, 'password': password},
      );

      final authResponse = AuthResponse.fromJson(response.data);

      if (authResponse.success && authResponse.data != null) {
        await StorageService.saveToken(authResponse.data!.token);
        await StorageService.saveUserId(authResponse.data!.user.id);
        await StorageService.saveUserData(authResponse.data!.user.toJson());
      }

      return authResponse;
    } on DioException catch (e) {
      if (e.response != null) {
        return AuthResponse.fromJson(e.response!.data);
      }
      return AuthResponse(
        success: false,
        error: ErrorData(
          code: 'NETWORK_ERROR',
          message: e.message ?? 'Network error occurred',
        ),
      );
    } catch (e) {
      return AuthResponse(
        success: false,
        error: ErrorData(code: 'UNKNOWN_ERROR', message: e.toString()),
      );
    }
  }

  Future<AuthResponse> register({
    required String username,
    required String email,
    required String password,
    required String displayName,
  }) async {
    try {
      final response = await _dio.post(
        ApiConstants.register,
        data: {
          'username': username,
          'email': email,
          'password': password,
          'display_name': displayName,
        },
      );

      final authResponse = AuthResponse.fromJson(response.data);

      if (authResponse.success && authResponse.data != null) {
        await StorageService.saveToken(authResponse.data!.token);
        await StorageService.saveUserId(authResponse.data!.user.id);
        await StorageService.saveUserData(authResponse.data!.user.toJson());
      }

      return authResponse;
    } on DioException catch (e) {
      if (e.response != null) {
        return AuthResponse.fromJson(e.response!.data);
      }
      return AuthResponse(
        success: false,
        error: ErrorData(
          code: 'NETWORK_ERROR',
          message: e.message ?? 'Network error occurred',
        ),
      );
    } catch (e) {
      return AuthResponse(
        success: false,
        error: ErrorData(code: 'UNKNOWN_ERROR', message: e.toString()),
      );
    }
  }

  Future<void> logout() async {
    await StorageService.clearAll();
  }

  Future<bool> isLoggedIn() async {
    final token = await StorageService.getToken();
    return token != null;
  }

  Future<UserModel?> getCurrentUser() async {
    final userData = await StorageService.getUserData();
    if (userData != null) {
      return UserModel.fromJson(userData);
    }
    return null;
  }
}
