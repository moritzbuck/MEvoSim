import unittest
from MEvoSim.classes.Gene import *

class TestGene(unittest.TestCase):

    def test_init(self):
        pass

class TestSimpleGene(unittest.TestCase):

    def test_Gene_init(self):
        gene = Gene()
        self.assertIsNotNone(gene._gene_id)

    def test_SimpleGene_init(self):
        gene = SimpleGene()
        self.assertIsNotNone(gene._gene_id)
        self.assertIn(gene._gene_id, SimpleGene._gene_dict)
        self.assertIs(SimpleGene._gene_dict[gene._gene_id], gene)
