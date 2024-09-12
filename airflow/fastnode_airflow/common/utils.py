import ast
import datetime
import time
import uuid
import json

from airflow.models import Variable
import fastnode_metrics


fastnode_status_config = fastnode_metrics.load_context('fastnode_status')


def is_production():
    return Variable.get('env', 'dev') == 'production'


def get_supported_languages():
    return fastnode_status_config['languages']


def get_unique_suffix():
    return '-{}-{}.json'.format(
        get_date_time_in_ISO(),
        uuid.uuid4().hex,
    )


def get_date_time_in_ISO():
    date_time = datetime.datetime.fromtimestamp(time.time())
    return date_time.isoformat() + 'Z'


def string_to_dict(string):
    return ast.literal_eval(string.replace('=', ':'))
