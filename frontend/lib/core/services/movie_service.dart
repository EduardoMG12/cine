import 'package:dio/dio.dart';
import 'api_service.dart';
import '../../features/movies/data/models/movie_model.dart';
import '../../features/movies/data/models/movie_detail_model.dart';

class MovieService {
  static final Dio _dio = ApiService.dio;

  static Future<List<MovieModel>> getTrendingMovies() async {
    try {
      final response = await _dio.get('/movies/trending');

      if (response.data['success'] == true) {
        final List<dynamic> data = response.data['data'];
        return data.map((json) => MovieModel.fromJson(json)).toList();
      } else {
        throw Exception('Failed to load trending movies');
      }
    } on DioException catch (e) {
      throw Exception('Failed to load trending movies: ${e.message}');
    }
  }

  static Future<MovieDetailModel> getMovieDetails(String imdbId) async {
    try {
      final response = await _dio.get('/omdb/$imdbId');
      return MovieDetailModel.fromJson(response.data);
    } on DioException catch (e) {
      throw Exception('Failed to load movie details: ${e.message}');
    }
  }

  static Future<List<MovieModel>> searchMovies(
    String query, {
    int page = 1,
  }) async {
    try {
      final response = await _dio.get(
        '/movies/search',
        queryParameters: {'q': query, 'page': page},
      );

      if (response.data['success'] == true) {
        final List<dynamic> data = response.data['data'];
        return data.map((json) => MovieModel.fromJson(json)).toList();
      } else {
        throw Exception('Failed to search movies');
      }
    } on DioException catch (e) {
      throw Exception('Failed to search movies: ${e.message}');
    }
  }
}
