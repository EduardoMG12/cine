import 'package:dio/dio.dart';
import 'api_service.dart';

class UserService {
  static final Dio _dio = ApiService.dio;

  // Get current user profile
  static Future<Map<String, dynamic>> getMe() async {
    try {
      print('👤 Getting current user profile...');
      final response = await _dio.get('/auth/me');
      print('👤 User profile response: ${response.data}');
      return response.data;
    } catch (e) {
      print('❌ Error getting user profile: $e');
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
      print('👤 Updating user profile...');

      final data = <String, dynamic>{};
      if (bio != null) data['bio'] = bio;
      if (displayName != null) data['display_name'] = displayName;
      if (isPrivate != null) data['is_private'] = isPrivate;
      if (profilePictureUrl != null)
        data['profile_picture_url'] = profilePictureUrl;
      if (theme != null) data['theme'] = theme;

      print('👤 Update data: $data');

      final response = await _dio.patch('/users/me', data: data);
      print('✅ Profile updated successfully: ${response.data}');
      return response.data;
    } catch (e) {
      print('❌ Error updating profile: $e');
      rethrow;
    }
  }
}
