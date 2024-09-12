import setuptools

setuptools.setup(
    name="fastnode-airflow-dags", # Replace with your own username
    version="0.0.1",
    author="Fastnode Team",
    description="Fastnode Airflow codes.",
    packages=setuptools.find_packages(),
    python_requires='>=3.6',
    include_package_data = True,

    entry_points = {
        'airflow.plugins': [
            'google_plugin = fastnode_airflow.plugins.google:GoogleSheetsPlugin'
        ]
    }
)
