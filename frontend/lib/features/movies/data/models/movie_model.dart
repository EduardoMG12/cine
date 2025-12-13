class MovieModel {
  final String id;
  final String externalApiId;
  final String title;
  final String overview;
  final String? posterUrl;
  final List<dynamic> genres;
  final bool adult;
  final DateTime createdAt;
  final DateTime updatedAt;

  MovieModel({
    required this.id,
    required this.externalApiId,
    required this.title,
    required this.overview,
    this.posterUrl,
    required this.genres,
    required this.adult,
    required this.createdAt,
    required this.updatedAt,
  });

  factory MovieModel.fromJson(Map<String, dynamic> json) {
    return MovieModel(
      id: json['id'] ?? '',
      externalApiId: json['external_api_id'] ?? '',
      title: json['title'] ?? '',
      overview: json['overview'] ?? '',
      posterUrl: json['poster_url'],
      genres: json['genres'] ?? [],
      adult: json['adult'] ?? false,
      createdAt: DateTime.parse(json['created_at']),
      updatedAt: DateTime.parse(json['updated_at']),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'external_api_id': externalApiId,
      'title': title,
      'overview': overview,
      'poster_url': posterUrl,
      'genres': genres,
      'adult': adult,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }
}
