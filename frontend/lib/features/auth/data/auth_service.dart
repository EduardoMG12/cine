import 'package:dio/dio.dart';
import '../../../core/constants/api_constants.dart';
import '../../../core/services/api_service.dart';
import '../../../core/services/storage_service.dart';
import '../domain/models/auth_response.dart';
import '../domain/models/user_model.dart';

class AuthService {
  final Dio _dio = ApiService.dio;

  // Login
  Future<AuthResponse> login({
    required String email,
    required String password,
  }) async {
    try {
      print('🔵 [AUTH_SERVICE] Iniciando login...');
      print('🔵 [AUTH_SERVICE] Email: $email');
      print('🔵 [AUTH_SERVICE] URL: ${ApiConstants.baseUrl}${ApiConstants.login}');
      
      final response = await _dio.post(
        ApiConstants.login,
        data: {
          'email': email,
          'password': password,
        },
      );

      print('✅ [AUTH_SERVICE] Login response recebida!');
      print('✅ [AUTH_SERVICE] Status: ${response.statusCode}');
      print('✅ [AUTH_SERVICE] Data: ${response.data}');

      final authResponse = AuthResponse.fromJson(response.data);
      
      if (authResponse.success && authResponse.data != null) {
        print('✅ [AUTH_SERVICE] Login bem-sucedido! Token salvo.');
        // Save token and user ID
        await StorageService.saveToken(authResponse.data!.token);
        await StorageService.saveUserId(authResponse.data!.user.id);
      } else {
        print('❌ [AUTH_SERVICE] Login falhou: ${authResponse.error?.message}');
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
        error: ErrorData(
          code: 'UNKNOWN_ERROR',
          message: e.toString(),
        ),
      );
    }
  }

  // Register
  Future<AuthResponse> register({
    required String username,
    required String email,
    required String password,
    required String displayName,
  }) async {
    try {
      print('🔵 [AUTH_SERVICE] Iniciando registro...');
      print('🔵 [AUTH_SERVICE] Username: $username');
      print('🔵 [AUTH_SERVICE] Email: $email');
      print('🔵 [AUTH_SERVICE] Display Name: $displayName');
      print('🔵 [AUTH_SERVICE] URL: ${ApiConstants.baseUrl}${ApiConstants.register}');
      
      final response = await _dio.post(
        ApiConstants.register,
        data: {
          'username': username,
          'email': email,
          'password': password,
          'display_name': displayName,
        },
      );

      print('✅ [AUTH_SERVICE] Register response recebida!');
      print('✅ [AUTH_SERVICE] Status: ${response.statusCode}');
      print('✅ [AUTH_SERVICE] Data: ${response.data}');

      final authResponse = AuthResponse.fromJson(response.data);
      
      if (authResponse.success && authResponse.data != null) {
        print('✅ [AUTH_SERVICE] Registro bem-sucedido! Token salvo.');
        // Save token and user ID
        await StorageService.saveToken(authResponse.data!.token);
        await StorageService.saveUserId(authResponse.data!.user.id);
      } else {
        print('❌ [AUTH_SERVICE] Registro falhou: ${authResponse.error?.message}');
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
        error: ErrorData(
          code: 'UNKNOWN_ERROR',
          message: e.toString(),
        ),
      );
    }
  }

  // Logout
  Future<void> logout() async {
    await StorageService.clearAll();
  }

  // Check if user is logged in
  Future<bool> isLoggedIn() async {
    final token = await StorageService.getToken();
    return token != null;
  }

  // Get current user (placeholder - would need a /me endpoint)
  Future<UserModel?> getCurrentUser() async {
    final userId = await StorageService.getUserId();
    return userId != null ? null : null; // TODO: Implement when /me endpoint is available
  }
}
