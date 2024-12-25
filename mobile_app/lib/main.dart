import 'package:flutter/material.dart';
import 'screens/documents_screen.dart';
import 'screens/register_screen.dart';
import 'screens/document_detail_screen.dart';
import 'screens/create_document_screen.dart';
import 'screens/login_screen.dart';
import 'screens/folder_screen.dart';
import 'screens/edit_document_screen.dart';
import 'screens/share_screen.dart';
import 'screens/shared_with_me_screen.dart';
import 'screens/shared_by_me_screen.dart';

void main() {
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Nexus App',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      initialRoute: '/login',
      routes: {
        '/login': (context) => LoginScreen(),
        '/register': (context) => RegisterScreen(),
        '/folders': (context) {
          final args = ModalRoute.of(context)?.settings.arguments as Map?;
          if (args == null || !args.containsKey('token')) {
            return const Scaffold(
              body: Center(child: Text('Error: token not provided')),
            );
          }
          return FoldersScreen(token: args['token']);
        },
        '/documents': (context) {
          final args = ModalRoute.of(context)?.settings.arguments as Map?;
          if (args == null ||
              !args.containsKey('token') ||
              !args.containsKey('folderId') ||
              !args.containsKey('folderName')) {
            return const Scaffold(
              body: Center(
                  child: Text('Ошибка: необходимые параметры не переданы')),
            );
          }
          return DocumentsScreen(
            token: args['token'],
            folderId: args['folderId'],
            folderName: args['folderName'],
          );
        },
        '/document-detail': (context) {
          final args = ModalRoute.of(context)?.settings.arguments as Map?;
          if (args == null ||
              !args.containsKey('id') ||
              !args.containsKey('token')) {
            return const Scaffold(
              body: Center(child: Text('Ошибка: аргументы не переданы')),
            );
          }
          return DocumentDetailScreen(id: args['id'], token: args['token']);
        },
        '/create-document': (context) {
          final args = ModalRoute.of(context)?.settings.arguments
              as Map<String, dynamic>?;
          if (args == null || !args.containsKey('token')) {
            return const Scaffold(
              body: Center(child: Text('Error: token not provided')),
            );
          }
          return CreateDocumentScreen(
            token: args['token'],
            folderId: args['folderId'],
            folderName: args['folderName'],
          );
        },
        '/edit-document': (context) {
          final args = ModalRoute.of(context)?.settings.arguments as Map?;
          if (args == null ||
              !args.containsKey('token') ||
              !args.containsKey('documentId') ||
              !args.containsKey('title') ||
              !args.containsKey('content')) {
            return const Scaffold(
              body:
                  Center(child: Text('Error: required arguments not provided')),
            );
          }
          return EditDocumentScreen(
            token: args['token'],
            documentId: args['documentId'],
            title: args['title'],
            content: args['content'],
          );
        },
        '/share': (context) {
          final args = ModalRoute.of(context)?.settings.arguments as Map?;
          if (args == null ||
              !args.containsKey('token') ||
              !args.containsKey('documentId') ||
              !args.containsKey('documentTitle')) {
            return const Scaffold(
              body:
                  Center(child: Text('Error: required arguments not provided')),
            );
          }
          return ShareScreen(
            token: args['token'],
            documentId: args['documentId'],
            documentTitle: args['documentTitle'],
          );
        },
        '/shared-with-me': (context) {
          final token = ModalRoute.of(context)?.settings.arguments as String?;
          if (token == null) {
            return const Scaffold(
              body: Center(child: Text('Error: token not provided')),
            );
          }
          return SharedWithMeScreen(token: token);
        },
        '/shared-by-me': (context) {
          final token = ModalRoute.of(context)?.settings.arguments as String?;
          if (token == null) {
            return const Scaffold(
              body: Center(child: Text('Error: token not provided')),
            );
          }
          return SharedByMeScreen(token: token);
        },
      },
      onUnknownRoute: (settings) {
        return MaterialPageRoute(
          builder: (context) => const Scaffold(
            body: Center(
              child: Text('Ошибка: маршрут не найден'),
            ),
          ),
        );
      },
    );
  }
}
