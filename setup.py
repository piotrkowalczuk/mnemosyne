from setuptools import setup
from subprocess import check_output

with open('VERSION.txt', 'r') as content_file:
    version = content_file.read()

    setup(
        name='mnemosyne-client',
        version=version[1:],
        description='mnemosyne service grpc client library',
        url='http://github.com/piotrkowalczuk/mnemosyne',
        author='Piotr Kowalczuk',
        author_email='p.kowalczuk.priv@gmail.com',
        license='MIT',
        packages=['mnemosynerpc'],
        install_requires=[
            'protobuf',
            'grpcio',
        ],
        zip_safe=False,
        keywords=['mnemosyne', 'grpc', 'session', 'service', 'client'],
        download_url='https://github.com/piotrkowalczuk/mnemosyne/archive/%s.tar.gz' % version.rstrip(),
      )