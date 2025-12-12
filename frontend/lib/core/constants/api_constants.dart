class ApiConstants {
  // Use 10.0.2.2 for Android emulator to access host machine
  // Use localhost for iOS simulator or web
  static const String baseUrl = 'http://10.0.2.2:8080/api/v1';

  // Auth endpoints
  static const String login = '/auth/login';
  static const String register = '/auth/register';
  static const String logout = '/auth/logout';
  static const String me = '/auth/me';
  static const String refreshToken = '/auth/refresh';

  // Movie endpoints
  static const String movies = '/movies';
  static String movieById(String id) => '/movies/$id';
  static const String searchMovies = '/movies/search';

  // User endpoints
  static const String profile = '/users/profile';
  static const String updateProfile = '/users/profile';

  // Social endpoints
  static const String friends = '/friends';
  static const String addFriend = '/friends/add';
  static const String removeFriend = '/friends/remove';

  // Lists endpoints
  static const String watchLater = '/lists/watch-later';
  static const String watched = '/lists/watched';
  static const String favorites = '/lists/favorites';

  // Match endpoints
  static const String createSession = '/match/session';
  static const String joinSession = '/match/session/join';
  static const String swipe = '/match/swipe';
}
