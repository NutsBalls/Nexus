import 'dart:convert';
import 'dart:io';
import 'package:http/http.dart' as http;
import 'package:http_parser/http_parser.dart';
import 'download_helper.dart' if (dart.library.html) 'download_helper_web.dart';

class ApiService {
  final String baseUrl = "http://backend:8080/api";

  Future<dynamic> post(String endpoint, Map<String, dynamic> data,
      [String? token]) async {
    final headers = {
      'Content-Type': 'application/json',
      if (token != null) 'Authorization': 'Bearer $token',
    };

    final response = await http.post(
      Uri.parse('$baseUrl/$endpoint'),
      body: jsonEncode(data),
      headers: headers,
    );

    if (response.statusCode == 200 || response.statusCode == 201) {
      return jsonDecode(response.body);
    } else {
      throw Exception('Error: ${response.statusCode}, ${response.body}');
    }
  }

  Future<dynamic> uploadFile(String endpoint, File file, String token) async {
    var uri = Uri.parse('$baseUrl/$endpoint');
    var request = http.MultipartRequest('POST', uri);

    request.headers['Authorization'] = 'Bearer $token';

    String fileExtension = file.path.split('.').last.toLowerCase();

    String contentType;
    switch (fileExtension) {
      case 'pdf':
        contentType = 'application/pdf';
        break;
      case 'doc':
      case 'docx':
        contentType = 'application/msword';
        break;
      case 'jpg':
      case 'jpeg':
        contentType = 'image/jpeg';
        break;
      case 'png':
        contentType = 'image/png';
        break;
      default:
        contentType = 'application/octet-stream';
    }

    request.files.add(
      await http.MultipartFile.fromPath(
        'file',
        file.path,
        contentType: MediaType.parse(contentType),
      ),
    );

    try {
      final streamedResponse = await request.send();
      final response = await http.Response.fromStream(streamedResponse);

      if (response.statusCode == 200 || response.statusCode == 201) {
        return jsonDecode(response.body);
      } else {
        throw Exception('Error: ${response.statusCode}, ${response.body}');
      }
    } catch (e) {
      throw Exception('Failed to upload file: $e');
    }
  }

  Future<dynamic> get(String endpoint, [String? token]) async {
    try {
      final headers = {
        'Content-Type': 'application/json',
        if (token != null) 'Authorization': 'Bearer $token',
      };

      final response = await http.get(
        Uri.parse('$baseUrl/$endpoint'),
        headers: headers,
      );

      if (response.statusCode == 200) {
        return jsonDecode(response.body);
      } else {
        throw Exception('Error: ${response.statusCode}, ${response.body}');
      }
    } catch (e) {
      throw Exception('Failed to fetch data: $e');
    }
  }

  Future<Map<String, dynamic>> uploadAttachment(
      File file, String documentId, String token) async {
    var uri = Uri.parse('$baseUrl/documents/$documentId/attachments');
    var request = http.MultipartRequest('POST', uri);

    request.headers['Authorization'] = 'Bearer $token';
    request.files.add(
      await http.MultipartFile.fromPath('file', file.path),
    );

    try {
      final streamedResponse = await request.send();
      final response = await http.Response.fromStream(streamedResponse);

      if (response.statusCode == 200 || response.statusCode == 201) {
        return jsonDecode(response.body);
      } else {
        throw Exception('Error: ${response.statusCode}, ${response.body}');
      }
    } catch (e) {
      throw Exception('Failed to upload attachment: $e');
    }
  }

  Future<List<Map<String, dynamic>>> getAttachments(
      String documentId, String token) async {
    try {
      final response = await get('documents/$documentId/attachments', token);
      return List<Map<String, dynamic>>.from(response);
    } catch (e) {
      throw Exception('Failed to get attachments: $e');
    }
  }

  Future<void> exportDocument(String documentId, String token) async {
    try {
      final response = await http.get(
        Uri.parse('$baseUrl/documents/$documentId/export'),
        headers: {
          'Authorization': 'Bearer $token',
        },
      );

      if (response.statusCode == 200) {
        final fileName = 'document_$documentId.json';
        await DownloadHelper.downloadFile(response.bodyBytes, fileName);
      } else {
        throw Exception('Error: ${response.statusCode}, ${response.body}');
      }
    } catch (e) {
      throw Exception('Failed to export document: $e');
    }
  }

  Future<Map<String, dynamic>> importDocument(File file, String token) async {
    var uri = Uri.parse('$baseUrl/documents/import');
    var request = http.MultipartRequest('POST', uri);

    request.headers['Authorization'] = 'Bearer $token';
    request.files.add(
      await http.MultipartFile.fromPath('document', file.path),
    );

    try {
      final streamedResponse = await request.send();
      final response = await http.Response.fromStream(streamedResponse);

      if (response.statusCode == 200) {
        return jsonDecode(response.body);
      } else {
        throw Exception('Error: ${response.statusCode}, ${response.body}');
      }
    } catch (e) {
      throw Exception('Failed to import document: $e');
    }
  }

  Future<void> downloadAttachment(String filePath, String token) async {
    try {
      final fileName = filePath.split('/').last;
      final encodedFileName = Uri.encodeComponent(fileName);

      final response = await http.get(
        Uri.parse('$baseUrl/download/$encodedFileName'),
        headers: {
          'Authorization': 'Bearer $token',
        },
      );

      if (response.statusCode == 200) {
        await DownloadHelper.downloadFile(response.bodyBytes, fileName);
      } else {
        throw Exception('Error: ${response.statusCode}, ${response.body}');
      }
    } catch (e) {
      throw Exception('Failed to download attachment: $e');
    }
  }

  Future<Map<String, dynamic>> uploadFileBytes(
    List<int> bytes,
    String fileName,
    String documentId,
    String token,
  ) async {
    var uri = Uri.parse('$baseUrl/documents/$documentId/attachments');
    var request = http.MultipartRequest('POST', uri);

    request.headers['Authorization'] = 'Bearer $token';

    request.files.add(
      http.MultipartFile.fromBytes(
        'file',
        bytes,
        filename: fileName,
      ),
    );

    try {
      final streamedResponse = await request.send();
      final response = await http.Response.fromStream(streamedResponse);

      if (response.statusCode == 200 || response.statusCode == 201) {
        return jsonDecode(response.body);
      } else {
        throw Exception('Error: ${response.statusCode}, ${response.body}');
      }
    } catch (e) {
      throw Exception('Failed to upload file: $e');
    }
  }

  Future<Map<String, dynamic>> getDocument(
      String documentId, String token) async {
    try {
      final response = await get('documents/$documentId', token);
      return response as Map<String, dynamic>;
    } catch (e) {
      throw Exception('Failed to get document: $e');
    }
  }

  Future<List<Map<String, dynamic>>> getDocumentAttachments(
    String documentId,
    String token,
  ) async {
    try {
      final response = await get('documents/$documentId/attachments', token);
      return List<Map<String, dynamic>>.from(response);
    } catch (e) {
      throw Exception('Failed to get attachments: $e');
    }
  }

  Future<void> deleteFolder(int folderId, String token) async {
    try {
      print('Deleting folder with ID: $folderId');
      final url = Uri.parse('$baseUrl/folders/$folderId');
      print('Delete URL: $url');

      final response = await http.delete(
        url,
        headers: {
          'Authorization': 'Bearer $token',
          'Content-Type': 'application/json',
        },
      );

      print('Response status: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode != 204 && response.statusCode != 200) {
        throw Exception('Error: ${response.statusCode}, ${response.body}');
      }
    } catch (e) {
      print('Delete error: $e');
      throw Exception('Failed to delete folder: $e');
    }
  }

  Future<void> deleteAttachment(int attachmentId, String token) async {
    try {
      final response = await http.delete(
        Uri.parse('$baseUrl/attachments/$attachmentId'),
        headers: {
          'Authorization': 'Bearer $token',
        },
      );

      if (response.statusCode != 204) {
        throw Exception('Error: ${response.statusCode}, ${response.body}');
      }
    } catch (e) {
      throw Exception('Failed to delete attachment: $e');
    }
  }

  Future<Map<String, dynamic>> createFolder(String name, String token) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/folders'),
        headers: {
          'Authorization': 'Bearer $token',
          'Content-Type': 'application/json',
        },
        body: jsonEncode({'name': name}),
      );

      if (response.statusCode == 201 || response.statusCode == 200) {
        return jsonDecode(response.body);
      } else {
        throw Exception('Error: ${response.statusCode}, ${response.body}');
      }
    } catch (e) {
      throw Exception('Failed to create folder: $e');
    }
  }

  Future<void> deleteDocument(int documentId, String token) async {
    try {
      final response = await http.delete(
        Uri.parse('$baseUrl/documents/$documentId'),
        headers: {
          'Authorization': 'Bearer $token',
          'Content-Type': 'application/json',
        },
      );

      if (response.statusCode != 204 && response.statusCode != 200) {
        throw Exception('Error: ${response.statusCode}, ${response.body}');
      }
    } catch (e) {
      throw Exception('Failed to delete document: $e');
    }
  }

  Future<List<dynamic>> getDocumentShares(int documentId, String token) async {
    try {
      final response = await http.get(
        Uri.parse('$baseUrl/documents/$documentId/shares'),
        headers: {
          'Authorization': 'Bearer $token',
          'Content-Type': 'application/json',
        },
      );

      if (response.statusCode == 200) {
        return jsonDecode(response.body);
      } else {
        throw Exception('Error: ${response.statusCode}, ${response.body}');
      }
    } catch (e) {
      throw Exception('Failed to get shares: $e');
    }
  }

  Future<void> shareDocument(
    int documentId,
    String email,
    String permission,
    String token,
  ) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/documents/$documentId/share'),
        headers: {
          'Authorization': 'Bearer $token',
          'Content-Type': 'application/json',
        },
        body: jsonEncode({
          'user_email': email,
          'permission': permission,
        }),
      );
      if (response.statusCode != 200) {
        throw Exception('Error: ${response.statusCode}, ${response.body}');
      }
    } catch (e) {
      throw Exception('Failed to share document: $e');
    }
  }

  Future<void> removeShare(int shareId, String token) async {
    try {
      final response = await http.delete(
        Uri.parse('$baseUrl/shares/$shareId'),
        headers: {
          'Authorization': 'Bearer $token',
          'Content-Type': 'application/json',
        },
      );

      if (response.statusCode == 204) {
        return;
      } else {
        throw Exception('Error: ${response.statusCode}, ${response.body}');
      }
    } catch (e) {
      throw Exception('Failed to remove share: $e');
    }
  }

  Future<Map<String, dynamic>> updateDocument(
    int documentId,
    Map<String, dynamic> data,
    String token,
  ) async {
    try {
      final response = await http.put(
        Uri.parse('$baseUrl/documents/$documentId'),
        headers: {
          'Authorization': 'Bearer $token',
          'Content-Type': 'application/json',
        },
        body: jsonEncode(data),
      );

      if (response.statusCode == 200) {
        return jsonDecode(response.body);
      } else {
        throw Exception('Error: ${response.statusCode}, ${response.body}');
      }
    } catch (e) {
      throw Exception('Failed to update document: $e');
    }
  }

  Future<Map<String, dynamic>> checkDocumentAccess(
    String documentId,
    String token,
  ) async {
    try {
      final response = await http.get(
        Uri.parse('$baseUrl/documents/$documentId/access'),
        headers: {
          'Authorization': 'Bearer $token',
          'Content-Type': 'application/json',
        },
      );

      if (response.statusCode == 200) {
        return Map<String, dynamic>.from(jsonDecode(response.body));
      } else {
        throw Exception('Error: ${response.statusCode}, ${response.body}');
      }
    } catch (e) {
      throw Exception('Failed to check document access: $e');
    }
  }

  Future<List<Map<String, dynamic>>> getSharedWithMe(String token) async {
    try {
      final response = await http.get(
        Uri.parse('$baseUrl/shares/shared-with-me'),
        headers: {
          'Authorization': 'Bearer $token',
          'Content-Type': 'application/json',
        },
      );

      if (response.statusCode == 200) {
        return List<Map<String, dynamic>>.from(jsonDecode(response.body));
      } else {
        throw Exception('Error: ${response.statusCode}, ${response.body}');
      }
    } catch (e) {
      throw Exception('Failed to fetch shared with me documents: $e');
    }
  }

  Future<List<Map<String, dynamic>>> getSharedByMe(String token) async {
    try {
      final response = await http.get(
        Uri.parse('$baseUrl/shares/shared-by-me'),
        headers: {
          'Authorization': 'Bearer $token',
          'Content-Type': 'application/json',
        },
      );

      if (response.statusCode == 200) {
        return List<Map<String, dynamic>>.from(jsonDecode(response.body));
      } else {
        throw Exception('Error: ${response.statusCode}, ${response.body}');
      }
    } catch (e) {
      throw Exception('Failed to fetch shared by me documents: $e');
    }
  }
}
