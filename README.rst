qtpm - Qt Package Manager Prototype
=======================================

Install
----------

.. code-block:: bash

   $ go get github.com/shibukawa/qtpm

Usage
----------

Create application template (with BSD LICENSE file)

.. code-block:: bash

   $ mkdir helloworld
   $ qtpm init app HelloWorld bsd
   $ ls
   CMakeExtra.txt   LICENSE.rst build       include     qtpackage.toml  src     test        vendor
   $ qtpm build 

Create library template (with MIT LICENSE file)

.. code-block:: bash

   $ mkdir awesomesdk
   $ cd awsomesdk
   $ qtpm init app AwesomeSDK mit
   CMakeExtra.txt   LICENSE.rst build       include     qtpackage.toml  src     test        vendor
   $ qtpm build 

Add files

.. code-block:: bash

   $ qtpm add class MyDialog@QDialog
   $ qtpm add test TestMyDialog

Qt Location
--------------

It uses CMake behind qtpm command to build. By default, Qt should be in default (``CMAKE_PREFIX_PATH``). If you put Qt out of the folder,
there are two ways to specify the Qt location.

1. qtpm sees environment variable ``QTDIR``:

   .. code-block:: bash

      $ QTDIR=~/Qt/5.5/clang_64 qtpm build

2. put ``qtpackage.user.toml`` that contains the following contents:

   .. code-block:: none

      qtdir = 'C:\Qt\5.5\mingw492_32'

If you don't use the both settings and Qt is not in ``CMAKE_PREFIX_PATH``, qtpm tries to search any locations.

Name Convention
--------------------

This tool behaves according to the convention over any configuration.

* Source and header files are under ``src`` folder.
* Tests are under ``test`` folder.
* Resources are under ``resource`` folder.
* One project folder includes one executable file or one shared library as the output.
* If there is ``src/main.cpp``, qtpm generates executable, otherwise shared library
* Each test classes are implemented in ``test/*_test.cpp`` files (no header files) and compiled into executable.
* Other ``.cpp`` files in ``test`` are treated test utility. They are linked with each test executables.

Project File
-----------------

Project file is written in TOML format.

* ``name``: Project name.
* ``author``: Author name.
* ``license``: License name.
* ``requires``: Dependency packages like ``'github.com/shibukawa/qtobubus'`` (this feature is not implemented yet).
* ``qtmodules``: Required qt modules like ``Widgets``, ``Xml``.
* ``version``: Version number like ``[1, 0, 0]``.

License
--------------

MIT
