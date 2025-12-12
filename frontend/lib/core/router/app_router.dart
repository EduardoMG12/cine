import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

// Import pages
import '../../features/auth/presentation/pages/login_page.dart';
import '../../features/auth/presentation/pages/register_page.dart';
import '../../features/auth/presentation/providers/auth_provider.dart';
import '../../features/home/presentation/pages/home_public_page.dart';
import '../../features/home/presentation/pages/home_private_page.dart';
import '../../features/movies/presentation/pages/movie_detail_page.dart';
import '../../features/profile/presentation/pages/profile_page.dart';
import '../../features/profile/presentation/pages/edit_profile_page.dart';
import '../../features/movies/presentation/pages/watch_later_page.dart';
import '../../features/movies/presentation/pages/watched_movies_page.dart';
import '../../features/social/presentation/pages/friends_page.dart';
import '../../features/match/presentation/pages/match_page.dart';

/// Router configuration
final routerProvider = Provider<GoRouter>((ref) {
  final authState = ref.watch(authStateProvider);

  return GoRouter(
    initialLocation: '/',
    redirect: (context, state) {
      final isLoggedIn = authState.isAuthenticated && authState.user != null;

      print(
        '🔄 [ROUTER] Redirect check - isLoggedIn: $isLoggedIn, user: ${authState.user?.username}, location: ${state.matchedLocation}',
      );

      final isAuthRoute = state.matchedLocation.startsWith('/auth');
      final isRootRoute = state.matchedLocation == '/';

      // If user is logged in and on root route, redirect to private home
      if (isLoggedIn && isRootRoute) {
        print(
          '🔄 [ROUTER] Usuário logado na rota raiz, redirecionando para /home',
        );
        return '/home';
      }

      // If user is logged in and on auth routes, redirect to home
      if (isLoggedIn && isAuthRoute) {
        print(
          '🔄 [ROUTER] Usuário logado em rota de auth, redirecionando para /home',
        );
        return '/home';
      }

      // If user is NOT logged in and trying to access protected routes, redirect to root
      if (!isLoggedIn && !isAuthRoute && state.matchedLocation != '/') {
        print(
          '🔄 [ROUTER] Usuário não logado tentando acessar rota protegida, redirecionando para /',
        );
        return '/';
      }

      // Allow navigation
      return null;
    },
    routes: [
      // Public routes
      GoRoute(path: '/', builder: (context, state) => const HomePublicPage()),

      GoRoute(
        path: '/movie/:id',
        builder: (context, state) {
          final movieId = state.pathParameters['id']!;
          return MovieDetailPage(movieId: movieId);
        },
      ),

      // Auth routes
      GoRoute(
        path: '/auth/login',
        builder: (context, state) => const LoginPage(),
      ),

      GoRoute(
        path: '/auth/register',
        builder: (context, state) => const RegisterPage(),
      ),

      // Private routes
      GoRoute(
        path: '/home',
        builder: (context, state) => const HomePrivatePage(),
      ),

      GoRoute(
        path: '/profile',
        builder: (context, state) => const ProfilePage(),
      ),

      GoRoute(
        path: '/profile/edit',
        builder: (context, state) => const EditProfilePage(),
      ),

      GoRoute(
        path: '/watch-later',
        builder: (context, state) => const WatchLaterPage(),
      ),

      GoRoute(
        path: '/watched',
        builder: (context, state) => const WatchedMoviesPage(),
      ),

      GoRoute(
        path: '/friends',
        builder: (context, state) => const FriendsPage(),
      ),

      GoRoute(path: '/match', builder: (context, state) => const MatchPage()),
    ],
    errorBuilder: (context, state) => Scaffold(
      body: Center(child: Text('Page not found: ${state.matchedLocation}')),
    ),
  );
});
