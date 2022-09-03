# frozen_string_literal: true
# require './helper'
require "haml"
module Graphql
  module Docs
    # This is a custom verion of Gitlab renderer for graphql-docs for Meshery.
    # Uses HAML templates to parse markdown and generate .md files.
    # It uses graphql-docs helpers and schema parser, more information in https://github.com/gjtorikian/graphql-docs.
    #
    # Arguments:
    #   schema - the GraphQL schema definition.
    #   output_dir: The folder where the markdown files will be saved
    #   template: The path of the haml template to be parsed
    class Renderer
      include Graphql::Docs::Helper

      attr_reader :schema

      def initialize(schema, output_dir, template)
        @output_dir = output_dir
        @template = template
        @schema = schema
        @layout = Haml::Engine.new(File.read(template))
        @parsed_schema = GraphQLDocs::Parser.new(schema, {}).parse
        @seen = Set.new
      end

      def contents
        # Render and remove an extra trailing new line
        @contents ||= @layout.render(self).sub!(/\n(?=\Z)/, '')
      end

      def write
        filename = File.join(@output_dir, 'graphql-api-reference.md')

        FileUtils.mkdir_p(@output_dir)
        File.write(filename, contents)
      end

      private

      def seen_type?(name)
        @seen.include?(name)
      end

      def seen_type!(name)
        @seen << name
      end
    end
  end
end
