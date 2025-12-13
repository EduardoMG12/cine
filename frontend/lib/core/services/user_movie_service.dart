import 'package:dio/dio.dart';
import 'api_service.dart';

class UserMovieService {
  static final Dio _dio = ApiService.dio;

  /// Toggle movie in favorites
  /// POST /api/v1/favorites
  static Future<Map<String, dynamic>> toggleFavoriteMovie(
    String movieId,
  ) async {
    try {
      final response = await _dio.post(
        '/favorites',
        data: {'movie_id': movieId},
      );

      return {
        'success': response.data['success'] ?? true,
        'message': response.data['message'] ?? 'Favorito atualizado',
        'added': response.data['added'] ?? false,
      };
    } on DioException catch (e) {
      print('Error toggling favorite: ${e.message}');
      throw Exception(
        'Erro ao atualizar favorito: ${e.response?.data['message'] ?? e.message}',
      );
    }
  }

  /// Toggle movie in watched list
  /// POST /api/v1/watched
  static Future<Map<String, dynamic>> toggleWatchedMovie(String movieId) async {
    try {
      final response = await _dio.post('/watched', data: {'movie_id': movieId});

      return {
        'success': response.data['success'] ?? true,
        'message': response.data['message'] ?? 'Status de assistido atualizado',
        'added': response.data['added'] ?? false,
      };
    } on DioException catch (e) {
      print('Error toggling watched: ${e.message}');
      throw Exception(
        'Erro ao atualizar assistido: ${e.response?.data['message'] ?? e.message}',
      );
    }
  }

  /// Get user's watched movies
  /// GET /api/v1/watched
  static Future<List<dynamic>> getWatchedMovies() async {
    try {
      final response = await _dio.get('/watched');
      print('📺 Watched movies response: ${response.data}');

      if (response.data['success'] == true) {
        final data = response.data['data'] ?? [];
        print('📺 Watched movies data count: ${data.length}');
        if (data.isNotEmpty) {
          print('📺 First movie: ${data[0]}');
        }
        return data;
      } else {
        throw Exception('Failed to load watched movies');
      }
    } on DioException catch (e) {
      print('❌ Error fetching watched movies: ${e.message}');
      print('❌ Response data: ${e.response?.data}');
      throw Exception('Erro ao carregar filmes assistidos: ${e.message}');
    }
  }

  /// Get user's favorite movies
  /// GET /api/v1/favorites
  static Future<List<dynamic>> getFavoriteMovies() async {
    try {
      final response = await _dio.get('/favorites');
      print('⭐ Favorite movies response: ${response.data}');

      if (response.data['success'] == true) {
        final data = response.data['data'] ?? [];
        print('⭐ Favorite movies data count: ${data.length}');
        if (data.isNotEmpty) {
          print('⭐ First movie: ${data[0]}');
        }
        return data;
      } else {
        throw Exception('Failed to load favorite movies');
      }
    } on DioException catch (e) {
      print('❌ Error fetching favorite movies: ${e.message}');
      print('❌ Response data: ${e.response?.data}');
      throw Exception('Erro ao carregar filmes favoritos: ${e.message}');
    }
  }
}
