import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import '../models/folder.dart';
import '../services/api_service.dart';

class FoldersScreen extends StatefulWidget {
  final String token;

  const FoldersScreen({Key? key, required this.token}) : super(key: key);

  @override
  _FoldersScreenState createState() => _FoldersScreenState();
}

class _FoldersScreenState extends State<FoldersScreen> {
  final ApiService _apiService = ApiService();
  List<Folder> folders = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadFolders();
  }

  Future<void> _loadFolders() async {
    setState(() => _isLoading = true);
    try {
      final response = await _apiService.get('folders', widget.token);
      setState(() {
        folders =
            (response as List).map((json) => Folder.fromJson(json)).toList();
        _isLoading = false;
      });
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to load folders: $e')),
        );
      }
      setState(() => _isLoading = false);
    }
  }

  Future<void> _createFolder() async {
    final nameController = TextEditingController();

    final name = await showDialog<String>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Create New Folder'),
        content: TextField(
          controller: nameController,
          decoration: const InputDecoration(
            labelText: 'Folder Name',
            border: OutlineInputBorder(),
          ),
          autofocus: true,
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () => Navigator.pop(context, nameController.text),
            child: const Text('Create'),
          ),
        ],
      ),
    );

    if (name != null && name.isNotEmpty) {
      try {
        final response = await _apiService.post(
          'folders',
          {'name': name},
          widget.token,
        );
        final newFolder = Folder.fromJson(response);
        setState(() {
          folders.add(newFolder);
        });
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Failed to create folder: $e')),
          );
        }
      }
    }
  }

  Future<void> _deleteFolder(int folderId) async {
    try {
      print('Starting folder deletion for ID: $folderId');
      await _apiService.deleteFolder(folderId, widget.token);

      setState(() {
        folders.removeWhere((folder) => folder.id == folderId);
      });

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Folder deleted successfully')),
        );
      }
    } catch (e) {
      print('Error during folder deletion: $e');
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to delete folder: $e')),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('My Folders'),
        actions: [
          IconButton(
            icon: const Icon(Icons.logout),
            onPressed: () => Navigator.pushReplacementNamed(context, '/login'),
          ),
          IconButton(
            icon: const Icon(Icons.share),
            onPressed: () {
              Navigator.pushNamed(
                context,
                '/shared-with-me',
                arguments: widget.token,
              );
            },
          ),
          IconButton(
            icon: const Icon(Icons.folder_shared),
            onPressed: () {
              Navigator.pushNamed(
                context,
                '/shared-by-me',
                arguments: widget.token,
              );
            },
          ),
        ],
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : folders.isEmpty
              ? Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      const Text('No folders yet'),
                      const SizedBox(height: 16),
                      ElevatedButton(
                        onPressed: _createFolder,
                        child: const Text('Create Folder'),
                      ),
                    ],
                  ),
                )
              : ListView.builder(
                  itemCount: folders.length,
                  itemBuilder: (context, index) {
                    final folder = folders[index];
                    return Dismissible(
                      key: Key(folder.id.toString()),
                      direction: DismissDirection.endToStart,
                      confirmDismiss: (direction) async {
                        return await showDialog<bool>(
                          context: context,
                          builder: (context) => AlertDialog(
                            title: const Text('Delete Folder'),
                            content: const Text(
                              'Are you sure you want to delete this folder? '
                              'All documents and files inside will be deleted permanently.',
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
                      onDismissed: (direction) {
                        _deleteFolder(folder.id);
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
                        leading: const Icon(Icons.folder),
                        title: Text(folder.name),
                        subtitle: Text(
                          'Created: ${DateFormat('MMM d, yyyy').format(folder.createdAt)}',
                        ),
                        trailing: _isLoading
                            ? const SizedBox(
                                width: 20,
                                height: 20,
                                child:
                                    CircularProgressIndicator(strokeWidth: 2),
                              )
                            : PopupMenuButton(
                                itemBuilder: (context) => [
                                  const PopupMenuItem(
                                    value: 'delete',
                                    child: Text('Delete'),
                                  ),
                                ],
                                onSelected: (value) async {
                                  if (value == 'delete') {
                                    final confirmed = await showDialog<bool>(
                                      context: context,
                                      builder: (context) => AlertDialog(
                                        title: const Text('Delete Folder'),
                                        content: const Text(
                                          'Are you sure you want to delete this folder? '
                                          'All documents and files inside will be deleted permanently.',
                                        ),
                                        actions: [
                                          TextButton(
                                            onPressed: () =>
                                                Navigator.pop(context, false),
                                            child: const Text('Cancel'),
                                          ),
                                          TextButton(
                                            onPressed: () =>
                                                Navigator.pop(context, true),
                                            style: TextButton.styleFrom(
                                              foregroundColor: Colors.red,
                                            ),
                                            child: const Text('Delete'),
                                          ),
                                        ],
                                      ),
                                    );

                                    if (confirmed == true) {
                                      await _deleteFolder(folder.id);
                                    }
                                  }
                                },
                              ),
                        onTap: _isLoading
                            ? null
                            : () {
                                Navigator.pushNamed(
                                  context,
                                  '/documents',
                                  arguments: {
                                    'token': widget.token,
                                    'folderId': folder.id,
                                    'folderName': folder.name,
                                  },
                                );
                              },
                      ),
                    );
                  },
                ),
      floatingActionButton: FloatingActionButton(
        onPressed: _createFolder,
        tooltip: 'Create Folder',
        child: const Icon(Icons.create_new_folder),
      ),
    );
  }
}
