import 'package:dio/dio.dart';
import 'api_service.dart';

class UserMovieService {
  static final Dio _dio = ApiService.dio;

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
      throw Exception(
        'Erro ao atualizar favorito: ${e.response?.data['message'] ?? e.message}',
      );
    }
  }

  static Future<Map<String, dynamic>> toggleWatchedMovie(String movieId) async {
    try {
      final response = await _dio.post('/watched', data: {'movie_id': movieId});

      return {
        'success': response.data['success'] ?? true,
        'message': response.data['message'] ?? 'Status de assistido atualizado',
        'added': response.data['added'] ?? false,
      };
    } on DioException catch (e) {
      throw Exception(
        'Erro ao atualizar assistido: ${e.response?.data['message'] ?? e.message}',
      );
    }
  }

  static Future<List<dynamic>> getWatchedMovies() async {
    try {
      final response = await _dio.get('/watched');

      if (response.data['success'] == true) {
        return response.data['data'] ?? [];
      } else {
        throw Exception('Failed to load watched movies');
      }
    } on DioException catch (e) {
      throw Exception('Erro ao carregar filmes assistidos: ${e.message}');
    }
  }

  static Future<List<dynamic>> getFavoriteMovies() async {
    try {
      final response = await _dio.get('/favorites');

      if (response.data['success'] == true) {
        return response.data['data'] ?? [];
      } else {
        throw Exception('Failed to load favorite movies');
      }
    } on DioException catch (e) {
      throw Exception('Erro ao carregar filmes favoritos: ${e.message}');
    }
  }
}
