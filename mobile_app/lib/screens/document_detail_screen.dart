import 'package:flutter/material.dart';
import 'package:mobile_app/screens/share_screen.dart';
import '../services/api_service.dart';
import 'package:file_picker/file_picker.dart';
import '../screens/edit_document_screen.dart';

class DocumentDetailScreen extends StatefulWidget {
  final String token;
  final int id;

  const DocumentDetailScreen({
    Key? key,
    required this.token,
    required this.id,
  }) : super(key: key);

  @override
  _DocumentDetailScreenState createState() => _DocumentDetailScreenState();
}

class _DocumentDetailScreenState extends State<DocumentDetailScreen> {
  final ApiService _apiService = ApiService();
  Map<String, dynamic>? document;
  List<Map<String, dynamic>> attachments = [];
  bool _isLoading = true;
  bool _isOwner = false;

  @override
  void initState() {
    super.initState();
    _loadDocument();
  }

  Future<void> _loadDocument() async {
    setState(() => _isLoading = true);
    try {
      final docResponse = await _apiService.get(
        'documents/${widget.id}',
        widget.token,
      );
      final attachResponse = await _apiService.get(
        'documents/${widget.id}/attachments',
        widget.token,
      );
      setState(() {
        document = docResponse;
        attachments = List<Map<String, dynamic>>.from(attachResponse);
        _isLoading = false;
      });
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Failed to load document: $e')),
      );
      setState(() => _isLoading = false);
    }
  }

  Future<void> _checkAccess() async {
    try {
      final accessResponse = await _apiService.checkDocumentAccess(
        widget.id.toString(),
        widget.token,
      );
      setState(() {
        _isOwner = accessResponse['isOwner'] ?? false;
      });
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Failed to check access: $e')),
      );
    }
  }

  Future<List<Map<String, dynamic>>> _getDocumentShares() async {
    try {
      final shares = await _apiService.getDocumentShares(
        widget.id,
        widget.token,
      );
      return List<Map<String, dynamic>>.from(shares);
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Failed to load shares: $e')),
      );
      return [];
    }
  }

  Future<void> _removeShare(int shareId) async {
    try {
      await _apiService.removeShare(shareId, widget.token);
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Share removed successfully')),
      );
      setState(() {});
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Failed to remove share: $e')),
      );
    }
  }

  Future<void> _pickAndUploadFile() async {
    try {
      FilePickerResult? result = await FilePicker.platform.pickFiles(
        allowMultiple: true,
        withData: true,
      );

      if (result != null) {
        setState(() => _isLoading = true);

        for (var file in result.files) {
          if (file.bytes != null) {
            await _apiService.uploadFileBytes(
              file.bytes!,
              file.name,
              widget.id.toString(),
              widget.token,
            );
          }
        }
        await _loadDocument();
      }
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Failed to upload file: $e')),
      );
      setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(document?['title'] ?? 'Document Details'),
        actions: [
          if (_isOwner)
            IconButton(
              icon: const Icon(Icons.edit),
              onPressed: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(
                    builder: (context) => EditDocumentScreen(
                      token: widget.token,
                      documentId: widget.id,
                      title: document?['title'] ?? '',
                      content: document?['content'] ?? '',
                    ),
                  ),
                );
              },
            ),
          if (_isOwner)
            IconButton(
              icon: const Icon(Icons.delete, color: Colors.red),
              onPressed: () async {
                final confirmed = await showDialog<bool>(
                  context: context,
                  builder: (context) => AlertDialog(
                    title: const Text('Delete Document'),
                    content: const Text(
                      'Are you sure you want to delete this document? '
                      'This action cannot be undone.',
                    ),
                    actions: [
                      TextButton(
                        onPressed: () => Navigator.pop(context, true),
                        style: TextButton.styleFrom(
                          foregroundColor: Colors.red,
                        ),
                        child: const Text('Delete'),
                      ),
                    ],
                  ),
                );

                if (confirmed == true) {
                  try {
                    await _apiService.deleteDocument(widget.id, widget.token);
                    if (mounted) {
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(
                            content: Text('Document deleted successfully')),
                      );
                      Navigator.pop(context);
                    }
                  } catch (e) {
                    if (mounted) {
                      ScaffoldMessenger.of(context).showSnackBar(
                        SnackBar(
                            content: Text('Failed to delete document: $e')),
                      );
                    }
                  }
                }
              },
            ),
          IconButton(
            icon: const Icon(Icons.share),
            onPressed: () {
              Navigator.push(
                context,
                MaterialPageRoute(
                  builder: (context) => ShareScreen(
                    token: widget.token,
                    documentId: widget.id,
                    documentTitle: document?['title'] ?? 'Document',
                  ),
                ),
              );
            },
          ),
        ],
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : SingleChildScrollView(
              padding: const EdgeInsets.all(16.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  if (!_isOwner) ...[
                    const SizedBox(height: 8),
                    const Text(
                      'You have read-only access to this document.',
                      style: TextStyle(
                        color: Colors.grey,
                        fontStyle: FontStyle.italic,
                      ),
                    ),
                  ],
                  Card(
                    child: Padding(
                      padding: const EdgeInsets.all(16.0),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            document?['title'] ?? 'Untitled',
                            style: Theme.of(context).textTheme.titleLarge,
                          ),
                          const SizedBox(height: 8),
                          Text(document?['content'] ?? 'No content'),
                        ],
                      ),
                    ),
                  ),
                  if (attachments.isNotEmpty) ...[
                    const SizedBox(height: 16),
                    Text(
                      'Attachments',
                      style: Theme.of(context).textTheme.titleLarge,
                    ),
                    const SizedBox(height: 8),
                    Card(
                      child: Column(
                        children: attachments.map((attachment) {
                          return ListTile(
                            leading: const Icon(Icons.attach_file),
                            title: Text(attachment['filename'] ?? ''),
                            subtitle: Text(
                              '${(attachment['size'] / 1024).toStringAsFixed(2)} KB',
                            ),
                            trailing: IconButton(
                              icon: const Icon(Icons.download),
                              onPressed: () async {
                                try {
                                  await _apiService.downloadAttachment(
                                    attachment['path'],
                                    widget.token,
                                  );
                                  ScaffoldMessenger.of(context).showSnackBar(
                                    const SnackBar(
                                        content: Text('Download started')),
                                  );
                                } catch (e) {
                                  ScaffoldMessenger.of(context).showSnackBar(
                                    SnackBar(
                                        content: Text(
                                            'Failed to download file: $e')),
                                  );
                                }
                              },
                            ),
                          );
                        }).toList(),
                      ),
                    ),
                  ],
                  if (_isOwner && document != null) ...[
                    const SizedBox(height: 16),
                    Text(
                      'Shared With',
                      style: Theme.of(context).textTheme.titleLarge,
                    ),
                    const SizedBox(height: 8),
                    FutureBuilder<List<Map<String, dynamic>>>(
                      future: _getDocumentShares(),
                      builder: (context, snapshot) {
                        if (snapshot.connectionState ==
                            ConnectionState.waiting) {
                          return const Center(
                              child: CircularProgressIndicator());
                        }
                        if (snapshot.hasError) {
                          return Center(
                              child: Text('Error: ${snapshot.error}'));
                        }
                        final shares = snapshot.data ?? [];
                        return Card(
                          child: Column(
                            children: shares.map((share) {
                              return ListTile(
                                title: Text(share['user']['username'] ??
                                    'Unknown User'),
                                subtitle:
                                    Text(share['user']['email'] ?? 'No email'),
                                trailing: IconButton(
                                  icon: const Icon(Icons.delete,
                                      color: Colors.red),
                                  onPressed: () => _removeShare(share['id']),
                                ),
                              );
                            }).toList(),
                          ),
                        );
                      },
                    ),
                  ],
                  const SizedBox(height: 16),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                    children: [
                      if (_isOwner)
                        ElevatedButton.icon(
                          onPressed: () async {
                            try {
                              await _apiService.exportDocument(
                                widget.id.toString(),
                                widget.token,
                              );
                              ScaffoldMessenger.of(context).showSnackBar(
                                const SnackBar(content: Text('Export started')),
                              );
                            } catch (e) {
                              ScaffoldMessenger.of(context).showSnackBar(
                                SnackBar(
                                    content:
                                        Text('Failed to export document: $e')),
                              );
                            }
                          },
                          icon: const Icon(Icons.download),
                          label: const Text('Export'),
                        ),
                      if (_isOwner)
                        ElevatedButton.icon(
                          onPressed: () => _pickAndUploadFile(),
                          icon: const Icon(Icons.attach_file),
                          label: const Text('Add Attachment'),
                        ),
                    ],
                  ),
                ],
              ),
            ),
    );
  }

  Widget _buildAttachmentsList() {
    return Card(
      child: Column(
        children: attachments.map((attachment) {
          return Dismissible(
            key: Key(attachment['id'].toString()),
            direction: DismissDirection.endToStart,
            confirmDismiss: (direction) async {
              return await showDialog<bool>(
                context: context,
                builder: (context) => AlertDialog(
                  title: const Text('Delete Attachment'),
                  content: const Text(
                    'Are you sure you want to delete this file? '
                    'This action cannot be undone.',
                  ),
                  actions: [
                    TextButton(
                      onPressed: () => Navigator.pop(context, false),
                      child: const Text('Cancel'),
                    ),
                    TextButton(
                      onPressed: () => Navigator.pop(context, true),
                      style: TextButton.styleFrom(
                        foregroundColor: Colors.red,
                      ),
                      child: const Text('Delete'),
                    ),
                  ],
                ),
              );
            },
            onDismissed: (direction) async {
              try {
                await _apiService.deleteAttachment(
                  attachment['id'],
                  widget.token,
                );
                _loadDocument();
                if (mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('File deleted successfully')),
                  );
                }
              } catch (e) {
                if (mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text('Failed to delete file: $e')),
                  );
                }
              }
            },
            background: Container(
              color: Colors.red,
              alignment: Alignment.centerRight,
              padding: const EdgeInsets.only(right: 16),
              child: const Icon(
                Icons.delete,
                color: Colors.white,
              ),
            ),
            child: ListTile(
              leading: const Icon(Icons.attach_file),
              title: Text(attachment['filename'] ?? ''),
              subtitle: Text(
                '${(attachment['size'] / 1024).toStringAsFixed(2)} KB',
              ),
              trailing: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  IconButton(
                    icon: const Icon(Icons.download),
                    onPressed: () async {
                      try {
                        await _apiService.downloadAttachment(
                          attachment['path'],
                          widget.token,
                        );
                        ScaffoldMessenger.of(context).showSnackBar(
                          const SnackBar(content: Text('Download started')),
                        );
                      } catch (e) {
                        ScaffoldMessenger.of(context).showSnackBar(
                          SnackBar(
                              content: Text('Failed to download file: $e')),
                        );
                      }
                    },
                  ),
                  if (_isOwner)
                    IconButton(
                      icon: const Icon(Icons.delete),
                      onPressed: () async {
                        final confirmed = await showDialog<bool>(
                          context: context,
                          builder: (context) => AlertDialog(
                            title: const Text('Delete File'),
                            content: const Text(
                              'Are you sure you want to delete this file? '
                              'This action cannot be undone.',
                            ),
                            actions: [
                              TextButton(
                                onPressed: () => Navigator.pop(context, false),
                                child: const Text('Cancel'),
                              ),
                              TextButton(
                                onPressed: () => Navigator.pop(context, true),
                                style: TextButton.styleFrom(
                                  foregroundColor: Colors.red,
                                ),
                                child: const Text('Delete'),
                              ),
                            ],
                          ),
                        );

                        if (confirmed == true) {
                          try {
                            await _apiService.deleteAttachment(
                              attachment['id'],
                              widget.token,
                            );
                            _loadDocument();
                            if (mounted) {
                              ScaffoldMessenger.of(context).showSnackBar(
                                const SnackBar(
                                    content: Text('File deleted successfully')),
                              );
                            }
                          } catch (e) {
                            if (mounted) {
                              ScaffoldMessenger.of(context).showSnackBar(
                                SnackBar(
                                    content: Text('Failed to delete file: $e')),
                              );
                            }
                          }
                        }
                      },
                    ),
                ],
              ),
            ),
          );
        }).toList(),
      ),
    );
  }
}
