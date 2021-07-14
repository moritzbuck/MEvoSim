import unittest
from MEvoSim.classes.Genome import *


class TestGenome(unittest.TestCase):

    def test_init(self):
        pass

class TestSimpleGenome(unittest.TestCase):

    def test_init(self):
        self.assertRaises(ValueError, SimpleGenome, -3)
