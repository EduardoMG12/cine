import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/services/movie_service.dart';
import '../../data/models/movie_detail_model.dart';

/// Provider para buscar detalhes de um filme específico
final movieDetailProvider = FutureProvider.family<MovieDetailModel, String>(
  (ref, imdbId) async {
    return await MovieService.getMovieDetails(imdbId);
  },
);
