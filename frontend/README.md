<div align="center">

# 📱 CineVerse Mobile App

### Flutter Mobile Application for Movie Enthusiasts

[![Flutter](https://img.shields.io/badge/Flutter-3.0+-02569B?style=for-the-badge&logo=flutter&logoColor=white)](https://flutter.dev/)
[![Dart](https://img.shields.io/badge/Dart-3.0+-0175C2?style=for-the-badge&logo=dart&logoColor=white)](https://dart.dev/)
[![Riverpod](https://img.shields.io/badge/Riverpod-2.0+-00A8E1?style=for-the-badge)](https://riverpod.dev/)
[![Android](https://img.shields.io/badge/Android-Ready-3DDC84?style=for-the-badge&logo=android&logoColor=white)](https://www.android.com/)
[![iOS](https://img.shields.io/badge/iOS-Ready-000000?style=for-the-badge&logo=apple&logoColor=white)](https://www.apple.com/ios/)

**Cross-Platform Mobile Client for CineVerse API**

[🇺🇸 English](#english) | [🇧🇷 Português](#português)

</div>

---

<a name="english"></a>

<details open>
<summary><h2>🇺🇸 ENGLISH VERSION</h2></summary>

<details>
<summary><h3>📋 Table of Contents</h3></summary>

- [About](#about-en)
- [Features](#features-en)
- [Architecture](#architecture-en)
- [Technologies](#technologies-en)
- [Project Structure](#structure-en)
- [Getting Started](#getting-started-en)
- [Development](#development-en)
- [Building](#building-en)
- [Testing](#testing-en)

</details>

<details>
<summary><h3>📖 About</h3></summary>

<a name="about-en"></a>

CineVerse Mobile is a cross-platform Flutter application that provides a seamless movie discovery and social experience. The app connects to the CineVerse API backend to deliver:

- 🎬 **Movie Discovery** - Browse trending movies and search by title
- ⭐ **Personal Lists** - Track favorite and watched movies
- 👤 **User Profile** - Manage account settings and preferences
- 🔐 **Secure Authentication** - JWT-based login and registration
- 🎨 **Modern UI** - Material Design 3 with smooth animations

**Key Highlights:**
- Clean Architecture with separation of concerns
- State management with Riverpod
- Offline-first with local caching
- Responsive design for all screen sizes
- Type-safe routing with GoRouter

</details>

<details>
<summary><h3>✨ Features</h3></summary>

<a name="features-en"></a>

#### Authentication
- ✅ User registration with email validation
- ✅ Secure login with JWT tokens
- ✅ Automatic token refresh
- ✅ Logout and session management

#### Movie Discovery
- ✅ Browse trending movies
- ✅ Search movies by title
- ✅ View detailed movie information
- ✅ Movie carousels (Trending, Popular, Recommended)
- ✅ Quick access to movie lists

#### User Lists
- ✅ Add/remove movies to favorites
- ✅ Track watched movies
- ✅ View favorites grid
- ✅ View watched movies grid
- ✅ Toggle movie status with one tap

#### User Profile
- ✅ View profile information
- ✅ Edit display name and bio
- ✅ Update profile picture
- ✅ Privacy settings (public/private profile)
- ✅ Theme preference (light/dark)
- ✅ Email verification status

#### UI/UX
- ✅ Material Design 3
- ✅ Smooth page transitions
- ✅ Pull-to-refresh
- ✅ Loading states
- ✅ Error handling with user feedback
- ✅ Responsive layout

</details>

<details>
<summary><h3>🏗️ Architecture</h3></summary>

<a name="architecture-en"></a>

#### Clean Architecture Layers

```
lib/
├── core/                    # Core functionality
│   ├── constants/          # API URLs, app constants
│   ├── router/             # GoRouter configuration
│   └── services/           # HTTP client, storage
│
├── features/               # Feature modules
│   ├── auth/              # Authentication
│   │   ├── data/          # Data sources & repositories
│   │   ├── domain/        # Models & entities
│   │   └── presentation/  # Pages, widgets, providers
│   │
│   ├── movies/            # Movie features
│   │   ├── data/          # Movie data layer
│   │   ├── domain/        # Movie models
│   │   └── presentation/  # Movie UI
│   │
│   ├── profile/           # User profile
│   │   └── presentation/  # Profile pages
│   │
│   └── home/              # Home & navigation
│       └── presentation/  # Home page, drawer
│
└── shared/                # Shared widgets
    └── widgets/           # Reusable components
```

#### Design Patterns

- **Provider Pattern** - State management with Riverpod
- **Repository Pattern** - Data access abstraction
- **Service Pattern** - Business logic encapsulation
- **Builder Pattern** - Widget composition
- **Observer Pattern** - Reactive state updates

#### State Management

```dart
// FutureProvider for async data
final movieProvider = FutureProvider<Movie>((ref) async {
  return await MovieService.getMovieById(id);
});

// StateProvider for simple state
final authStateProvider = StateNotifierProvider<AuthNotifier, AuthState>(
  (ref) => AuthNotifier(),
);
```

</details>

<details>
<summary><h3>🛠️ Technologies</h3></summary>

<a name="technologies-en"></a>

#### Core Framework
- **Flutter 3.0+** - Cross-platform UI framework
- **Dart 3.0+** - Programming language

#### State Management
- **flutter_riverpod ^2.4.0** - State management solution
- **riverpod_annotation** - Code generation for providers

#### Networking
- **dio ^5.4.0** - HTTP client
- **retrofit** - Type-safe REST client (optional)

#### Navigation
- **go_router ^13.0.0** - Declarative routing
- **flutter_native_splash** - Splash screen

#### Storage
- **flutter_secure_storage ^9.0.0** - Secure credential storage
- **shared_preferences** - Local key-value storage

#### UI Components
- **cached_network_image** - Image caching
- **flutter_svg** - SVG support
- **shimmer** - Loading placeholders

#### Development Tools
- **build_runner** - Code generation
- **flutter_launcher_icons** - App icon generation
- **flutter_lints** - Linting rules

</details>

<details>
<summary><h3>📁 Project Structure</h3></summary>

<a name="structure-en"></a>

```
frontend/
├── android/                # Android platform code
├── ios/                    # iOS platform code
├── lib/
│   ├── core/
│   │   ├── constants/
│   │   │   └── api_constants.dart
│   │   ├── router/
│   │   │   └── app_router.dart
│   │   └── services/
│   │       ├── api_service.dart
│   │       ├── storage_service.dart
│   │       ├── movie_service.dart
│   │       ├── user_service.dart
│   │       └── user_movie_service.dart
│   │
│   ├── features/
│   │   ├── auth/
│   │   │   ├── data/
│   │   │   │   └── auth_service.dart
│   │   │   ├── domain/
│   │   │   │   └── models/
│   │   │   │       ├── user_model.dart
│   │   │   │       └── auth_response.dart
│   │   │   └── presentation/
│   │   │       ├── pages/
│   │   │       │   ├── login_page.dart
│   │   │       │   └── register_page.dart
│   │   │       └── providers/
│   │   │           └── auth_provider.dart
│   │   │
│   │   ├── movies/
│   │   │   ├── data/
│   │   │   │   └── models/
│   │   │   │       ├── movie_model.dart
│   │   │   │       └── movie_detail_model.dart
│   │   │   └── presentation/
│   │   │       └── pages/
│   │   │           ├── movie_detail_page.dart
│   │   │           ├── watch_later_page.dart
│   │   │           └── watched_movies_page.dart
│   │   │
│   │   ├── profile/
│   │   │   └── presentation/
│   │   │       └── pages/
│   │   │           └── profile_page.dart
│   │   │
│   │   └── home/
│   │       └── presentation/
│   │           └── pages/
│   │               ├── home_public_page.dart
│   │               └── home_private_page.dart
│   │
│   ├── shared/
│   │   └── widgets/
│   │       └── movie_card.dart
│   │
│   └── main.dart
│
├── test/                   # Unit & widget tests
├── pubspec.yaml           # Dependencies
└── analysis_options.yaml  # Linting rules
```

</details>

<details>
<summary><h3>🚀 Getting Started</h3></summary>

<a name="getting-started-en"></a>

#### Prerequisites

- Flutter SDK 3.0 or higher
- Dart SDK 3.0 or higher
- Android Studio / Xcode (for platform-specific builds)
- CineVerse API running (backend)

#### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/EduardoMG12/cine.git
   cd cine/frontend
   ```

2. **Install dependencies**
   ```bash
   flutter pub get
   ```

3. **Configure API endpoint**
   
   Edit `lib/core/constants/api_constants.dart`:
   ```dart
   class ApiConstants {
     static const String baseUrl = 'http://YOUR_API_HOST:8080/api/v1';
   }
   ```

4. **Run the app**
   ```bash
   # List available devices
   flutter devices
   
   # Run on specific device
   flutter run -d <device-id>
   
   # Run in debug mode
   flutter run --debug
   
   # Run in release mode
   flutter run --release
   ```

#### Environment Setup

**For Android Development:**
- Install Android Studio
- Set up Android SDK
- Create AVD (Android Virtual Device)
- Enable USB debugging on physical device

**For iOS Development (macOS only):**
- Install Xcode
- Set up iOS Simulator
- Configure code signing
- Trust developer certificate on device

</details>

<details>
<summary><h3>💻 Development</h3></summary>

<a name="development-en"></a>

#### Running in Development Mode

```bash
# Hot reload enabled
flutter run

# With specific device
flutter run -d chrome          # Web
flutter run -d emulator-5554   # Android emulator
flutter run -d "iPhone 14"     # iOS simulator
```

#### Code Generation

```bash
# Generate code for providers, models, etc.
flutter pub run build_runner build

# Watch mode (regenerate on file changes)
flutter pub run build_runner watch

# Delete conflicting outputs
flutter pub run build_runner build --delete-conflicting-outputs
```

#### Linting & Formatting

```bash
# Analyze code
flutter analyze

# Format code
flutter format .

# Fix auto-fixable issues
dart fix --apply
```

#### Clean Build

```bash
# Clean build artifacts
flutter clean

# Reinstall dependencies
flutter pub get

# Full rebuild
flutter run
```

#### Debugging

- Use **Flutter DevTools** for debugging
- Access at: `http://localhost:9100`
- Features:
  - Widget inspector
  - Performance profiler
  - Network inspector
  - Logging console

</details>

<details>
<summary><h3>📦 Building</h3></summary>

<a name="building-en"></a>

#### Android Build

```bash
# Build APK (debug)
flutter build apk --debug

# Build APK (release)
flutter build apk --release

# Build App Bundle (for Play Store)
flutter build appbundle --release

# Build split APKs per ABI
flutter build apk --split-per-abi
```

**Output locations:**
- APK: `build/app/outputs/flutter-apk/app-release.apk`
- AAB: `build/app/outputs/bundle/release/app-release.aab`

#### iOS Build

```bash
# Build iOS app (release)
flutter build ios --release

# Build IPA for distribution
flutter build ipa --release
```

**Output locations:**
- iOS app: `build/ios/iphoneos/Runner.app`
- IPA: `build/ios/ipa/`

#### Web Build

```bash
# Build for web
flutter build web --release

# Build with web renderer
flutter build web --web-renderer canvaskit  # Better graphics
flutter build web --web-renderer html      # Smaller size
```

**Output location:** `build/web/`

</details>

<details>
<summary><h3>🧪 Testing</h3></summary>

<a name="testing-en"></a>

#### Run Tests

```bash
# Run all tests
flutter test

# Run with coverage
flutter test --coverage

# Run specific test file
flutter test test/auth/auth_service_test.dart

# Run tests in watch mode
flutter test --watch
```

#### Test Structure

```
test/
├── unit/              # Unit tests
│   ├── services/
│   └── models/
├── widget/            # Widget tests
│   └── pages/
└── integration/       # Integration tests
    └── flows/
```

#### Writing Tests

```dart
// Unit test example
void main() {
  group('MovieService', () {
    test('should fetch trending movies', () async {
      final movies = await MovieService.getTrendingMovies();
      expect(movies, isNotEmpty);
    });
  });
}

// Widget test example
testWidgets('LoginPage displays correctly', (tester) async {
  await tester.pumpWidget(MyApp());
  expect(find.text('Login'), findsOneWidget);
});
```

</details>

</details>

---

<a name="português"></a>

<details>
<summary><h2>🇧🇷 VERSÃO EM PORTUGUÊS</h2></summary>

<details>
<summary><h3>📋 Índice</h3></summary>

- [Sobre](#about-pt)
- [Funcionalidades](#features-pt)
- [Arquitetura](#architecture-pt)
- [Tecnologias](#technologies-pt)
- [Estrutura do Projeto](#structure-pt)
- [Começando](#getting-started-pt)
- [Desenvolvimento](#development-pt)
- [Build](#building-pt)
- [Testes](#testing-pt)

</details>

<details>
<summary><h3>📖 Sobre</h3></summary>

<a name="about-pt"></a>

CineVerse Mobile é uma aplicação Flutter multiplataforma que oferece uma experiência perfeita de descoberta de filmes e interação social. O app se conecta ao backend da API CineVerse para entregar:

- 🎬 **Descoberta de Filmes** - Navegue filmes em alta e pesquise por título
- ⭐ **Listas Pessoais** - Acompanhe filmes favoritos e assistidos
- 👤 **Perfil de Usuário** - Gerencie configurações e preferências da conta
- 🔐 **Autenticação Segura** - Login e registro baseados em JWT
- 🎨 **UI Moderna** - Material Design 3 com animações suaves

**Destaques:**
- Clean Architecture com separação de responsabilidades
- Gerenciamento de estado com Riverpod
- Offline-first com cache local
- Design responsivo para todos os tamanhos de tela
- Roteamento type-safe com GoRouter

</details>

<details>
<summary><h3>✨ Funcionalidades</h3></summary>

<a name="features-pt"></a>

#### Autenticação
- ✅ Registro de usuário com validação de email
- ✅ Login seguro com tokens JWT
- ✅ Atualização automática de token
- ✅ Logout e gerenciamento de sessão

#### Descoberta de Filmes
- ✅ Navegue filmes em alta
- ✅ Pesquise filmes por título
- ✅ Visualize informações detalhadas do filme
- ✅ Carrosséis de filmes (Em Alta, Populares, Recomendados)
- ✅ Acesso rápido às listas de filmes

#### Listas de Usuário
- ✅ Adicionar/remover filmes dos favoritos
- ✅ Rastrear filmes assistidos
- ✅ Ver grid de favoritos
- ✅ Ver grid de filmes assistidos
- ✅ Alternar status do filme com um toque

#### Perfil de Usuário
- ✅ Ver informações do perfil
- ✅ Editar nome de exibição e bio
- ✅ Atualizar foto de perfil
- ✅ Configurações de privacidade (perfil público/privado)
- ✅ Preferência de tema (claro/escuro)
- ✅ Status de verificação de email

#### UI/UX
- ✅ Material Design 3
- ✅ Transições de página suaves
- ✅ Pull-to-refresh
- ✅ Estados de carregamento
- ✅ Tratamento de erros com feedback ao usuário
- ✅ Layout responsivo

</details>

<details>
<summary><h3>🏗️ Arquitetura</h3></summary>

<a name="architecture-pt"></a>

#### Camadas da Clean Architecture

```
lib/
├── core/                    # Funcionalidade principal
│   ├── constants/          # URLs da API, constantes
│   ├── router/             # Configuração GoRouter
│   └── services/           # Cliente HTTP, storage
│
├── features/               # Módulos de funcionalidades
│   ├── auth/              # Autenticação
│   │   ├── data/          # Fontes de dados & repositórios
│   │   ├── domain/        # Modelos & entidades
│   │   └── presentation/  # Páginas, widgets, providers
│   │
│   ├── movies/            # Funcionalidades de filmes
│   │   ├── data/          # Camada de dados de filmes
│   │   ├── domain/        # Modelos de filmes
│   │   └── presentation/  # UI de filmes
│   │
│   ├── profile/           # Perfil de usuário
│   │   └── presentation/  # Páginas de perfil
│   │
│   └── home/              # Home & navegação
│       └── presentation/  # Página inicial, drawer
│
└── shared/                # Widgets compartilhados
    └── widgets/           # Componentes reutilizáveis
```

#### Padrões de Design

- **Padrão Provider** - Gerenciamento de estado com Riverpod
- **Padrão Repository** - Abstração de acesso a dados
- **Padrão Service** - Encapsulamento de lógica de negócio
- **Padrão Builder** - Composição de widgets
- **Padrão Observer** - Atualizações de estado reativas

</details>

<details>
<summary><h3>🛠️ Tecnologias</h3></summary>

<a name="technologies-pt"></a>

#### Framework Principal
- **Flutter 3.0+** - Framework de UI multiplataforma
- **Dart 3.0+** - Linguagem de programação

#### Gerenciamento de Estado
- **flutter_riverpod ^2.4.0** - Solução de gerenciamento de estado
- **riverpod_annotation** - Geração de código para providers

#### Rede
- **dio ^5.4.0** - Cliente HTTP
- **retrofit** - Cliente REST type-safe (opcional)

#### Navegação
- **go_router ^13.0.0** - Roteamento declarativo
- **flutter_native_splash** - Tela de splash

#### Armazenamento
- **flutter_secure_storage ^9.0.0** - Armazenamento seguro de credenciais
- **shared_preferences** - Armazenamento local chave-valor

#### Componentes UI
- **cached_network_image** - Cache de imagens
- **flutter_svg** - Suporte SVG
- **shimmer** - Placeholders de carregamento

#### Ferramentas de Desenvolvimento
- **build_runner** - Geração de código
- **flutter_launcher_icons** - Geração de ícones do app
- **flutter_lints** - Regras de linting

</details>

<details>
<summary><h3>🚀 Começando</h3></summary>

<a name="getting-started-pt"></a>

#### Pré-requisitos

- Flutter SDK 3.0 ou superior
- Dart SDK 3.0 ou superior
- Android Studio / Xcode (para builds específicos de plataforma)
- API CineVerse rodando (backend)

#### Instalação

1. **Clone o repositório**
   ```bash
   git clone https://github.com/EduardoMG12/cine.git
   cd cine/frontend
   ```

2. **Instale as dependências**
   ```bash
   flutter pub get
   ```

3. **Configure o endpoint da API**
   
   Edite `lib/core/constants/api_constants.dart`:
   ```dart
   class ApiConstants {
     static const String baseUrl = 'http://SEU_HOST_API:8080/api/v1';
   }
   ```

4. **Execute o app**
   ```bash
   # Liste dispositivos disponíveis
   flutter devices
   
   # Execute em dispositivo específico
   flutter run -d <device-id>
   
   # Execute em modo debug
   flutter run --debug
   
   # Execute em modo release
   flutter run --release
   ```

</details>

<details>
<summary><h3>💻 Desenvolvimento</h3></summary>

<a name="development-pt"></a>

#### Executando em Modo de Desenvolvimento

```bash
# Hot reload habilitado
flutter run

# Com dispositivo específico
flutter run -d chrome          # Web
flutter run -d emulator-5554   # Emulador Android
flutter run -d "iPhone 14"     # Simulador iOS
```

#### Geração de Código

```bash
# Gerar código para providers, models, etc.
flutter pub run build_runner build

# Modo watch (regerar em mudanças de arquivo)
flutter pub run build_runner watch

# Deletar saídas conflitantes
flutter pub run build_runner build --delete-conflicting-outputs
```

#### Linting & Formatação

```bash
# Analisar código
flutter analyze

# Formatar código
flutter format .

# Corrigir problemas auto-corrigíveis
dart fix --apply
```

#### Build Limpo

```bash
# Limpar artefatos de build
flutter clean

# Reinstalar dependências
flutter pub get

# Rebuild completo
flutter run
```

</details>

<details>
<summary><h3>📦 Build</h3></summary>

<a name="building-pt"></a>

#### Build Android

```bash
# Build APK (debug)
flutter build apk --debug

# Build APK (release)
flutter build apk --release

# Build App Bundle (para Play Store)
flutter build appbundle --release

# Build APKs divididos por ABI
flutter build apk --split-per-abi
```

**Localizações de saída:**
- APK: `build/app/outputs/flutter-apk/app-release.apk`
- AAB: `build/app/outputs/bundle/release/app-release.aab`

#### Build iOS

```bash
# Build app iOS (release)
flutter build ios --release

# Build IPA para distribuição
flutter build ipa --release
```

**Localizações de saída:**
- App iOS: `build/ios/iphoneos/Runner.app`
- IPA: `build/ios/ipa/`

</details>

<details>
<summary><h3>🧪 Testes</h3></summary>

<a name="testing-pt"></a>

#### Executar Testes

```bash
# Executar todos os testes
flutter test

# Executar com cobertura
flutter test --coverage

# Executar arquivo de teste específico
flutter test test/auth/auth_service_test.dart

# Executar testes em modo watch
flutter test --watch
```

#### Estrutura de Testes

```
test/
├── unit/              # Testes unitários
│   ├── services/
│   └── models/
├── widget/            # Testes de widget
│   └── pages/
└── integration/       # Testes de integração
    └── flows/
```

</details>

</details>

---

<div align="center">

### 📄 License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.

### 👥 Team

Developed with ❤️ by the CineVerse Team

**Federal Institute of Paraná - Campus Palmas**

---

**[⬆ Back to top](#-cineverse-mobile-app)**

</div>

