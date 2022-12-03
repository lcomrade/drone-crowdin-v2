This Drone plugin has nothing to do with [jonasfranz/drone-crowdin](https://github.com/jonasfranz/drone-crowdin).
The plugin from Jonasfranz uses an outdated version of the Crowdin v1 API,
so it does not work correctly with new Crowdin projects and API keys.

This Drone plugin (`drone-crowdin-v2`) works with the current version of Crowdin API v2.
The source code of this software is also written from scratch
and is available under the [AGPLv3 license](LICENSE).


### Upload files to Crowdin
```yaml
kind: pipeline
name: translate

trigger:
  branch:
    - main
  event:
    - push

steps:
  - name: upload
    pull: always
    image: git.lcomrade.su/root/drone-crowdin-v2
    settings:
      crowdin_key:
        from_secret: crowdin_key

      project_id: 553341
      # You can use the project name if you don't know the ID:
      #   project_name: Lenpaste

      target: upload

      upload_files: {"internal/web/data/locale/en.locale": "en.ini"}
      # 1. Format: {"LOCAL_FILE_PATH", "CROWDIN_FILE_NAME"}
      # 2. If the file exists in Crowdin, a new revision will be created.
      # 3. Upload multiple files:
      #   upload_files: {"internal/web/data/locale/en.locale": "en.ini", "internal/web/data/locale/ru.locale": "ru.ini"}
```

INFO: You must create a secret `crowdin_key` and put the API token [obtained from Crowdin](https://crowdin.com/settings#api-key) into it.


### Download translate from Crowdin and push it to Git
```yaml
kind: pipeline
name: translate

trigger:
  branch:
    - main
  event:
    - push

steps:
  - name: download
    pull: always
    image: git.lcomrade.su/root/drone-crowdin-v2
    settings:
      crowdin_key:
        from_secret: crowdin_key

      project_id: 553341
      # You can use the project name if you don't know the ID:
      #   project_name: Lenpaste

      target: download
      download_to: internal/web/data/locale/

      # Extra settings:
      #   download_skip_untranslated_strings: false
      #   download_skip_untranslated_files: false
      #   download_export_approved_only: false


  - name: push
    pull: always
    image: appleboy/drone-git-push
    settings:
      author_email: "translatebot@example.org"
      author_name: Translate [Bot]
      branch: main
      commit: true
      commit_message: "[skip ci] Updated translations"
      remote: "git@example.org:root/lenpaste.git"
      ssh_key:
        from_secret: ci_ssh_key
```

INFO: You must create a secret `crowdin_key` and put the API token [obtained from Crowdin](https://crowdin.com/settings#api-key) into it.

INFO: Read more about the appleboy/drone-git-push plugin [here](https://github.com/appleboy/drone-git-push).
