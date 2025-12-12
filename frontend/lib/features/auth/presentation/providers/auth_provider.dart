import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../data/auth_service.dart';
import '../../domain/models/user_model.dart';

// Auth service provider
final authServiceProvider = Provider<AuthService>((ref) => AuthService());

// Auth state provider - usando NotifierProvider ao invés de StateNotifierProvider
final authStateProvider = NotifierProvider<AuthNotifier, AuthState>(() {
  return AuthNotifier();
});

// Auth state class
class AuthState {
  final UserModel? user;
  final bool isLoading;
  final String? errorMessage;

  AuthState({
    this.user,
    this.isLoading = false,
    this.errorMessage,
  });

  bool get isAuthenticated => user != null;

  AuthState copyWith({
    UserModel? user,
    bool? isLoading,
    String? errorMessage,
  }) {
    return AuthState(
      user: user ?? this.user,
      isLoading: isLoading ?? this.isLoading,
      errorMessage: errorMessage ?? this.errorMessage,
    );
  }
}

// Auth notifier - usando Notifier ao invés de StateNotifier
class AuthNotifier extends Notifier<AuthState> {
  AuthService get _authService => ref.read(authServiceProvider);

  @override
  AuthState build() {
    _checkAuthStatus();
    return AuthState();
  }

  // Check if user is already logged in
  Future<void> _checkAuthStatus() async {
    final isLoggedIn = await _authService.isLoggedIn();
    if (isLoggedIn) {
      // TODO: Fetch user data when /me endpoint is available
      // For now, just mark as authenticated with null user
      state = state.copyWith(user: null);
    }
  }

  // Login
  Future<bool> login(String email, String password) async {
    print('🟡 [AUTH_PROVIDER] Login iniciado');
    state = state.copyWith(isLoading: true, errorMessage: null);

    try {
      print('🟡 [AUTH_PROVIDER] Chamando auth service...');
      final response = await _authService.login(
        email: email,
        password: password,
      );

      print('🟡 [AUTH_PROVIDER] Response recebida: success=${response.success}');

      if (response.success && response.data != null) {
        print('✅ [AUTH_PROVIDER] Login bem-sucedido! Usuário: ${response.data!.user.username}');
        state = state.copyWith(
          user: response.data!.user,
          isLoading: false,
        );
        return true;
      } else {
        print('❌ [AUTH_PROVIDER] Login falhou: ${response.error?.message}');
        state = state.copyWith(
          isLoading: false,
          errorMessage: response.error?.message ?? 'Login failed',
        );
        return false;
      }
    } catch (e) {
      print('❌ [AUTH_PROVIDER] Exceção no login: $e');
      state = state.copyWith(
        isLoading: false,
        errorMessage: e.toString(),
      );
      return false;
    }
  }

  // Register
  Future<bool> register({
    required String username,
    required String email,
    required String password,
    required String displayName,
  }) async {
    print('🟡 [AUTH_PROVIDER] Register iniciado');
    state = state.copyWith(isLoading: true, errorMessage: null);

    try {
      final response = await _authService.register(
        username: username,
        email: email,
        password: password,
        displayName: displayName,
      );

      if (response.success && response.data != null) {
        state = state.copyWith(
          user: response.data!.user,
          isLoading: false,
        );
        return true;
      } else {
        state = state.copyWith(
          isLoading: false,
          errorMessage: response.error?.message ?? 'Registration failed',
        );
        return false;
      }
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        errorMessage: e.toString(),
      );
      return false;
    }
  }

  // Logout
  Future<void> logout() async {
    await _authService.logout();
    state = AuthState();
  }

  // Clear error
  void clearError() {
    state = state.copyWith(errorMessage: null);
  }
}

