Terraform Plugin for Beanstalk
==============================

This repository contains an out-of-tree Terraform plugin that allows the
creation and management of repositories, users and teams in the private Git
hosting solution [Beanstalk](http://beanstalkapp.com/).

This provides a means to manage basic Beanstalk configuration as code under
version control, rather than configuring it manually via the web UI.

Usage Example
-------------

Here's how to create a repository and a user and then create a team that grants
the user access to the repository:

```
provider "beanstalk" {
    account_name = "example"
}

resource "beanstalk_repository" "example" {
    name = "example"
    title = "Example Repository"
    color_label = "red"
}

resource "beanstalk_user" "example" {
    username = "example"
    email = "example@example.com"
    name = "Example User"
}

resource "beanstalk_team" "example" {
    name = "Example Team"
    user_ids = [
        "${beanstalk_user.example.id}"
    ]

    repository_permissions {
        repository_id = "${beanstalk_repository.example.id}"
        can_write = true
    }
}
```

See later in this document for more detailed documentation about each resource.

Installation
------------

Just like Terraform itself, this plugin is a Go application and it requires
the same minimal version of Go as Terraform, which is 1.4 at the time of
writing.

You'll need to set up a Go development environment as documented in
[How to Write Go Code](https://golang.org/doc/code.html).

Then you can install this plugin:

* ``go install github.com/saymedia/terraform-beanstalk/terraform-provider-beanstalk``

If successful then this should create a program ``terraform-provider-beanstalk``
in ``$GOPATH/bin``. You can copy this program anywhere you like to install
it, and then create a ``.terraformrc`` in your home directory that tells
Terraform where to find it:

```
providers {
    beanstalk = "/path/to/terraform-provider-beanstalk"
}
```

After this Terraform should automatically discover the plugin and the plugin's
resources should become available for use in Terraform configurations.

Usage
-----

The plugin is a provider which requires some configuration. The provider's
declaration gives the name of your Beanstalk account:

```
provider "beanstalk" {
    account_name = "example"
}
```

In addition to this you must also set the ``BEANSTALK_USERNAME`` and
``BEANSTALK_ACCESS_TOKEN`` environment variables to provide your Beanstalk
authentication credentials. You can get your access token from your Beanstalk
user settings.

Contributing
------------

This plugin was created primarily for internal use and was released only in the
hope that it would be useful to others. Code contributions are welcomed, but we
do not expect to add further features to this plugin unless they are needed by
our internal teams.

This codebase currently has no tests of any kind so we would appreciate any
manual testing that contributors can do before submitting contributions.

Please ensure that contributions are idiomatic Go and are formatting using
``gofmt``.

License
-------

The MIT License (MIT)

Copyright (c) 2015 Say Media Inc

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

Available Resources
===================

Repository
----------

The ``beanstalk_repository`` resource allows Beanstalk repositories to be
created. It supports the following parameters:

* ``name`` (required): the name of the repository that will be used in URLs
  referring to the repository. This must be unique across all repositories
  in your account.

* ``title`` (required): the name of the repository that will be shown in the
  Beanstalk web UI.

* ``color_label`` (optional): the color to use to represent the repository
  in the Beanstalk web UI. A limited number of colors are available as
  described in [the Beanstalk API docs](http://api.beanstalkapp.com/repository.html).
  Defaults to "white".

* ``default_git_branch`` (optional): the git branch name to use as the default
  branch, which will be checked out by default when users clone the repository
  and offered as the default target branch for code review requests.
  Defaults to "master". Only used when ``vcs`` is set to "git".

* ``vcs`` (optional): Either "git" or "subversion" depending on what kind of
  repository is desired. Defaults to "git".

* ``create_svn_structure`` (optional): Boolean defining whether to create the
  usual "trunk", "branches" and "tags" directory structure in a Subversion
  repository. Defaults to ``false``. Has no effect when ``vcs`` is not set to
  "subversion".

Once created, repository resources export the following attributes:

* ``id``: the id of the repository in Beanstalk.

* ``url``: the URL at which the created repository can be found. This is the
  URL that can be provided to either ``git clone`` or ``svn checkout``,
  depending on which VCS was chosen.

*Beanstalk does not allow repositories to be deleted via the API*. In order to
delete a repository that is managed by Beanstalk, log in to the web UI and
delete it from there, and then run ``terraform refresh`` to allow Terraform
to notice that the repository is delete. You can then remove the resource from
your Terraform configuration and run ``terraform apply`` to make it stick.

User
----

The ``beanstalk_user`` resource allows users to be invited to Beanstalk and
their accounts to then be managed by Terraform. It supports the following
parameters:

* ``username`` (required): The username that the user will use to log in. This
  must be unique within your account.

* ``name`` (required): The full name of the user that will be displayed in the
  Beanstalk UI. Beanstalk requires this to be at least two words separated by
  a space.

* ``email`` (required): The email address of the user. This must be unique
  within your account.

* ``account_admin`` (optional): Boolean defining whether the user will have
  administrative access to the Beanstalk account.

* ``timezone`` (optional): The name of the timezone that will be used to show
  this user times within the Beanstalk UI. This should be set to one of the
  strings from the timezone drop-down within the Beanstalk profile editing UI.
  It defaults to "London".

Once created, user resources export the following attributes:

* ``id``: the id of the user in Beanstalk.

* ``first_name``: Beanstalk's idea of the user's first name, extracted from
  the ``name`` parameter.

* ``last_name``: Beanstalk's idea of the user's last name, extracted from
  the ``name`` parameter.

Setting a user's password via Terraform is not supported, since users should
select their own passwords. When a new user resource is created,
*an invitation will be sent to the provided email address* and the user will
then have the opportunity to sign up, set a password, and (if enabled on your
Beanstalk account) configure two-factor authentication.

The user will be able to make adjustments to some of the settings that
Terraform controls via the Beanstalk UI. Terraform will reset these back to
the configured values next time it is run. Users should be advised not to
change these values in the UI but instead to change them (or have them changed)
in the Terraform configuration.

Team
----

The ``beanstalk_team`` resource allows Beanstalk teams to be
created. It supports the following parameters:

* ``name`` (required): The name of the team to be displayed in the Terraform UI.
  This must be unique within your account.

* ``user_ids`` (required): array of ids of users that will be members of the
  team.

* ``repository_permissions`` (required): a nested configuration block describing
  this team's permissions on a particular repository. Described in more detail
  below.

* ``color_label`` (optional): as with ``color_label`` on repositories, the
  color to use for the team in the Beanstalk UI.

The ``repository_permissions`` block has the following sub-parameters:

* ``repository_id`` (required): The id of the repository that this set of
  permissions applies to. Must be unique within the permissions of a
  particular team.

* ``can_write`` (optional): Boolean defining whether the team members have
  write access to the repository. Defaults to ``false``.

* ``can_deploy`` (optional): Boolean defining whether the team members have
  access to run deployments for the repository. Defaults to ``false``.

* ``can_configure_deployments`` (optional): Boolean defining whether the team
  members have access to configure deployments for the repository.
  Defaults to ``false``.

Once created, user resources export the following attribute:

* ``id``: the id of the team in Beanstalk.

Repository Code Review Settings
-------------------------------

The ``beanstalk_repository_code_review_settings`` resource allows the code
review settings for a Beanstalk repository to be managed in Terraform.
It supports the following parameters:

* ``repository_id`` (required): The id of the repository whose settings are
  managed by this block. This must be unique within a given Terraform
  configuration, since each repository has only one set of code review
  settings.

* ``unanimous_approval`` (optional): Boolean defining whether unanimous
  approval of all assigned reviewers is required before a review can be
  accepted. If ``false``, only one approval is required. Defaults to ``false``.

* ``auto_reopen`` (optional): Booleaning defining whether code review requests
  are automatically reopened when new commentary is added. Defaults to
  ``false``.

* ``default_assignee_user_ids`` (optional): Array of ids of Beanstalk users
  who will be automatically assigned to any new code review request.

* ``default_watching_user_ids`` (optional): Array of ids of Beanstalk users
  who will be automatically added as watchers on any new code review request.

* ``default_watching_team_ids`` (optional): Array of ids of Beanstalk teams
  whose members will be automatically added as watchers on any new code review
  request.

Creating more than one ``beanstalk_repository_code_review_settings`` resource
for the same repository is not supported and will result in a configuration
that is unable to converge, since all resources will separately try to control
the same underlying data.

