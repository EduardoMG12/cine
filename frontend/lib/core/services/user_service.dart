import 'package:dio/dio.dart';
import 'api_service.dart';

class UserService {
  static final Dio _dio = ApiService.dio;

  // Get current user profile
  static Future<Map<String, dynamic>> getMe() async {
    try {
      final response = await _dio.get('/auth/me');
      return response.data;
    } catch (e) {
      rethrow;
    }
  }

  // Update user profile
  static Future<Map<String, dynamic>> updateProfile({
    String? bio,
    String? displayName,
    bool? isPrivate,
    String? profilePictureUrl,
    String? theme,
  }) async {
    try {
      final data = <String, dynamic>{};
      if (bio != null) data['bio'] = bio;
      if (displayName != null) data['display_name'] = displayName;
      if (isPrivate != null) data['is_private'] = isPrivate;
      if (profilePictureUrl != null)
        data['profile_picture_url'] = profilePictureUrl;
      if (theme != null) data['theme'] = theme;

      final response = await _dio.patch('/users/me', data: data);
      return response.data;
    } catch (e) {
      rethrow;
    }
  }
}
