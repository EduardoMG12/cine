import 'dart:convert';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class StorageService {
  static const _storage = FlutterSecureStorage();

  static const String _tokenKey = 'auth_token';
  static const String _userIdKey = 'user_id';
  static const String _userDataKey = 'user_data';

  // Token operations
  static Future<void> saveToken(String token) async {
    await _storage.write(key: _tokenKey, value: token);
  }

  static Future<String?> getToken() async {
    return await _storage.read(key: _tokenKey);
  }

  static Future<void> deleteToken() async {
    await _storage.delete(key: _tokenKey);
  }

  // User ID operations
  static Future<void> saveUserId(String userId) async {
    await _storage.write(key: _userIdKey, value: userId);
  }

  static Future<String?> getUserId() async {
    return await _storage.read(key: _userIdKey);
  }

  static Future<void> deleteUserId() async {
    await _storage.delete(key: _userIdKey);
  }

  // User data operations
  static Future<void> saveUserData(Map<String, dynamic> userData) async {
    await _storage.write(key: _userDataKey, value: jsonEncode(userData));
  }

  static Future<Map<String, dynamic>?> getUserData() async {
    final data = await _storage.read(key: _userDataKey);
    if (data != null) {
      return jsonDecode(data) as Map<String, dynamic>;
    }
    return null;
  }

  static Future<void> deleteUserData() async {
    await _storage.delete(key: _userDataKey);
  }

  // Clear all auth data
  static Future<void> clearAll() async {
    await _storage.deleteAll();
  }
}
